/*
@Description:
@File : model
@Author : blue_sky_12138
@Version: 1.0.0
@Date : 2021/11/29 21:54
*/

package AGID

type Collection struct {
	Info  Info  `json:"info"`
	Items Items `json:"item"`
}

type Info struct {
	ID     string `json:"_postman_id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type Item struct {
	Name     string     `json:"name"`
	Items    Items      `json:"item"` // 文件夹套文件夹
	Req      Request    `json:"request"`
	Examples []Response `json:"response"`
}

type Items []*Item

type Request struct {
	Method  string `json:"method"`
	Headers []KVTD `json:"headers"`
	URL     URL    `json:"url"`
	// 因里面的结构体字段不完全相同，故用map存储
	// 其中mode字段存储实际发送的请求体的类型
	// 以mode字段的值为字段名的，存储了请求参数键值对
	// options字段存储预设的其他请求的设置
	Body map[string]interface{} `json:"body"`
	// 因里面的结构体字段不完全相同，故用map存储
	// 其中type字段存储实际发送的鉴权的类型
	// 以type字段的值为字段名的，存储了请求参数键值对
	// [未验证]options字段存储预设的其他请求的设置
	Auth map[string]interface{} `json:"auth"`
}

type KVTD struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type URL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Path     []string `json:"path"`
	Query    []KVTD   `json:"query"`
}

type Response struct {
	Name            string  `json:"name"`
	OriReq          Request `json:"originalRequest"`
	Status          string  `json:"status"`
	Code            int     `json:"code"`
	PreviewLanguage string  `json:"_postman_previewlanguage"`
	Headers         []KVTD
	Body            string `json:"body"` // 拿到的body是转义后的
}
