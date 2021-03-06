/*
@Description:
@File : logic
@Author : blue_sky_12138
@Version: 1.0.0
@Date : 2021/12/7 10:55
*/

package AGID

import (
	"AGID/markdown"
	"AGID/qjson"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
)

const PathSeparator = string(os.PathSeparator)

const schemaPostmanV21 = `https://schema.getpostman.com/json/collection/v2.1.0/collection.json`

var (
	ErrSchemaNoSuppose = errors.New("schema no suppose")
)

func LoadAndSaveDic(path string, recursion ...bool) (err error) {
	ss, err := LoadCollectionDic(path, recursion...)
	if err != nil {
		return
	}

	cols, err := AnalyseRaws(ss)
	if err != nil {
		return
	}

	for _, col := range cols {
		file, err := os.Create(`./` + col.Name() + ".md")
		if err != nil {
			return err
		}

		_, err = file.WriteString(col.String())
		if err != nil {
			return err
		}

		file.Close()
	}

	return nil
}

func LoadAndSave(filePath string) (err error) {
	s, err := LoadCollectionFile(filePath)
	if err != nil {
		return
	}

	col, err := AnalyseRaw(s)
	if err != nil {
		return
	}

	file, err := os.Create(`./` + col.Name() + ".md")
	if err != nil {
		return
	}

	_, err = file.WriteString(col.String())
	if err != nil {
		return
	}

	file.Close()

	return nil
}

func LoadCollectionDic(path string, recursion ...bool) ([]string, error) {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	rec := false

	if len(recursion) > 0 {
		rec = recursion[0]
	}

	res := make([]string, 0)
	for _, fir := range dir {
		if fir.IsDir() {
			if rec {
				cts, err := LoadCollectionDic(path)
				if err != nil {
					return nil, err
				}
				res = append(res, cts...)
			}
			continue
		}

		ct, err := LoadCollectionFile(path + PathSeparator + fir.Name())
		if err != nil {
			return nil, err
		}
		res = append(res, ct)
	}
	return res, nil
}

func LoadCollectionFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

type AGIDParser interface {
	Parse(raw string) (AGIDCollection, error)
}

type parse struct{}

func (p *parse) Parse(raw string) (c AGIDCollection, err error) {
	switch {
	case strings.ContainsAny(raw, schemaPostmanV21):
		var col Collection
		err = qjson.UnmarshalString(raw, &col)
		c = &col

	default:
		return nil, ErrSchemaNoSuppose
	}

	return
}

func AnalyseRaws(raws []string, p ...AGIDParser) ([]AGIDCollection, error) {
	res := make([]AGIDCollection, 0, len(raws))
	for _, raw := range raws {
		mk, err := AnalyseRaw(raw, p...)
		if err != nil {
			return nil, err
		}
		res = append(res, mk)
	}
	return res, nil
}

func AnalyseRaw(raw string, p ...AGIDParser) (res AGIDCollection, err error) {
	if len(p) > 0 {
		return p[0].Parse(raw)
	}

	return (&parse{}).Parse(raw)
}

type AGIDCollection interface {
	Name() string
	String() string
}

func (c *Collection) Name() string {
	return c.Info.Name
}

func (c *Collection) String() string {
	b := strings.Builder{}

	writeWithEnter(&b, markdown.H1(c.Name()))

	writeWithEnter(&b, c.Items.String(2))

	return b.String()
}

func (i *Item) markdownName(deep int) (s string) {
	switch deep {
	case 1:
		s = markdown.H1(i.Name)

	case 2:
		s = markdown.H2(i.Name)

	case 3:
		s = markdown.H3(i.Name)

	case 4:
		s = markdown.H4(i.Name)

	case 5:
		s = markdown.H5(i.Name)

	case 6:
		s = markdown.H6(i.Name)

	default:
		s = markdown.Bold(i.Name)
	}

	return
}

func (i Items) String(deep int) string {
	b := strings.Builder{}
	for _, item := range i {
		writeWithNewLine(&b, item.String(deep))
	}

	return b.String()
}

func (i *Item) String(deep int) string {
	b := strings.Builder{}

	// ?????????
	writeWithEnter(&b, i.markdownName(deep))

	// ???????????????????????????
	if i.Items != nil {
		writeWithNewLine(&b, i.Items.String(deep+1))

		return b.String()
	}

	// ????????????
	writeWithEnter(&b, markdown.Bold("???????????????"))
	writeWithEnter(&b, i.reqPath())

	// ?????????
	writeWithEnter(&b, markdown.Bold("????????????"))
	writeWithEnter(&b, i.reqHeader())

	// ????????????
	writeWithEnter(&b, markdown.Bold("???????????????"))
	writeWithEnter(&b, i.reqParam())

	// ????????????????????????????????????????????????
	// ??????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????????
	if i.haveResponse() {
		// ????????????
		if re := i.reqExample(); re != "" {
			writeWithEnter(&b, markdown.Bold("???????????????"))
			writeWithEnter(&b, re)
		}

		// ????????????
		if rp := i.respParam(); rp != "" {
			writeWithEnter(&b, markdown.Bold("???????????????"))
			writeWithEnter(&b, rp)
		}

		// ????????????
		writeWithEnter(&b, markdown.Bold("???????????????"))
		writeWithEnter(&b, i.respExample())
	}

	b.WriteByte('\n')
	return b.String()
}

func (i *Item) reqPath() (res string) {
	return i.Req.path()
}

func (r *Request) path() (res string) {
	res = r.Method + " /" + path.Join(r.URL.Path...)
	res = markdown.CodeBlock(res, "http")
	return
}

func (i *Item) reqHeader() string {
	return i.Req.header()
}

func (r *Request) header() string {
	ft := newDefaultFormTable()
	ft.SetText(0, 2, "??????")

	if r.Auth != nil {
		ft.AddLineAndSet("Authorization", "???", "Bearer $token", "??????token")
	}

	for _, header := range r.Headers {
		ft.AddLineAndSet(header.Key, "???", "/", "")
	}

	if ft.Line <= 1 {
		return "???"
	}

	return ft.String()
}

func (i *Item) reqParam() string {
	return i.Req.param()
}

func (r *Request) param() string {
	if s := r.parseQuery(); s != "" {
		return markdown.Bold("query??????") + "\n" + s
	}

	if s := r.parsePath(); s != "" {
		return markdown.Bold("path??????") + "\n" + s
	}

	return r.parseMode()
}

func (r *Request) parseQuery() string {
	if len(r.URL.Query) <= 0 {
		return ""
	}

	ft := newDefaultFormTable()
	for _, kvtd := range r.URL.Query {
		ft.AddLineAndSet(kvtd.Key, "???", kvtd.Type, kvtd.Description)
	}

	return ft.String()
}

func (r *Request) parsePath() string {
	var haveParam bool

	ft := newDefaultFormTable()
	for _, value := range r.URL.Path {
		if len(value) > 2 && value[0] == '{' && value[len(value)-1] == '}' {
			ft.AddLineAndSet(value[1:len(value)-1], "???", "", "")
			haveParam = true
		}
	}

	if !haveParam {
		return ""
	}

	return ft.String()
}

func (r *Request) parseMode() string {
	mode, modeExist := r.Body["mode"]
	if !modeExist {
		return "???"
	}

	modeString := mode.(string)
	var paramString string
	switch modeString {
	case "raw":
		paramString = "json??????"

	case "urlencoded":
		paramString = "x-www-form-urlencoded"

	case "formdata":
		paramString = "form??????"

	default:
		paramString = modeString + "??????"
	}
	paramString = markdown.Bold(paramString) + "\n"

	ft := newDefaultFormTable()

	if modeValue, exist := r.Body[modeString]; exist {
		switch values := modeValue.(type) {
		case string:
			m := parseJson(values)
			for k, v := range m {
				typ := getObjectType(reflect.ValueOf(v).Type())
				ft.AddLineAndSet(k, "???", typ, "")
			}

		case []interface{}:
			array := parseArray(values)
			for _, v := range array {
				ft.AddLineAndSet(v.Key, "???", v.Type, v.Description)
			}
		}

	}

	if ft.Line <= 1 {
		return "???"
	}

	return paramString + ft.String()
}

func parseJson(raw string) map[string]interface{} {
	// ??????postman??????????????????
	reg, _ := regexp.Compile(`({{.*?}})`)
	values := reg.ReplaceAllString(raw, "0")

	// ??????
	m := make(map[string]interface{})
	qjson.UnmarshalString(values, &m)

	return m
}

func getObjectType(typ reflect.Type) (typName string) {
	switch typ.Kind() {
	case reflect.Slice:
		typName = getObjectType(typ.Elem()) + "??????"

	case reflect.Array:
		typName = getObjectType(typ.Elem()) + "??????"

	case reflect.Map:
		typName = "??????????????????"

	case reflect.Interface:
		typName = "??????????????????"

	default:
		typName = typ.Name()
	}

	return
}

func parseArray(data []interface{}) []KVTD {
	res := make([]KVTD, 0, len(data))
	for _, dt := range data {
		m := dt.(map[string]interface{})

		k, _ := m["key"].(string)
		v, _ := m["value"].(string)
		t, _ := m["type"].(string)
		d, _ := m["description"].(string)
		res = append(res, KVTD{
			Key:         k,
			Value:       v,
			Type:        t,
			Description: d,
		})
	}

	return res
}

func (i *Item) haveResponse() bool {
	return len(i.Examples) > 0
}

func (i *Item) reqExample() string {
	if !i.haveResponse() {
		return ""
	}

	resp := i.getSuccessResp()
	if resp.OriReq.Body == nil {
		return ""
	}

	mode, modeExist := i.Req.Body["mode"]
	if !modeExist {
		return ""
	}

	modeString := mode.(string)
	if modeString != "raw" {
		return ""
	}

	return markdown.CodeBlock(i.Req.Body[modeString].(string), "json")
}

func (i *Item) respParam() string {
	if !i.haveResponse() {
		return ""
	}

	m := i.getSuccessResp().param()
	ft := newDefaultFormTable()

	for k, v := range m {
		typ := getObjectType(reflect.ValueOf(v).Type())
		ft.AddLineAndSet(k, "???", typ, "")
	}

	if ft.Line > 1 {
		return ft.String()
	} else {
		return "???"
	}
}

func (r Response) param() (m map[string]interface{}) {
	rawM := parseJson(r.Body)
	i, exist := rawM["data"]
	if !exist {
		return
	}

	m = make(map[string]interface{})
	switch value := i.(type) {
	case map[string]interface{}:
		for k, v := range value {
			m[k] = v
		}

	default:
		m["data(????????????????????????)"] = value
	}

	return
}

func (i *Item) getSuccessResp() Response {
	if !i.haveResponse() {
		return Response{}
	}

	if len(i.Examples) == 1 {
		return i.Examples[0]
	}

	rs := make([]Response, 0, len(i.Examples))
	for _, example := range i.Examples {
		if example.Code != http.StatusOK {
			continue
		}

		if strings.Contains(example.Name, "??????") || strings.Contains(example.Name, "??????") {
			return example
		}

		rs = append(rs, example)
	}

	if len(rs) == 0 {
		rs = append(rs, i.Examples[0])
	}

	return rs[0]
}

func (i *Item) respExample() string {
	if !i.haveResponse() {
		return ""
	}

	b := strings.Builder{}
	for _, example := range i.Examples {
		writeWithEnter(&b, example.Name)
		writeWithEnter(&b, markdown.CodeBlock(example.Body, example.PreviewLanguage))
	}

	return b.String()
}
