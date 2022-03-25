/*
@Description:
@File : markdown
@Author : blue_sky_12138
@Version: 1.0.0
@Date : 2021/12/7 11:24
*/

package markdown

import (
	"strconv"
	"strings"
)

func H1(text string) string {
	return "# " + text
}

func H2(text string) string {
	return "## " + text
}

func H3(text string) string {
	return "### " + text
}

func H4(text string) string {
	return "#### " + text
}

func H5(text string) string {
	return "##### " + text
}

func H6(text string) string {
	return "###### " + text
}

func Quote(text string) string {
	return "> " + text
}

func CodeBlock(text string, language ...string) string {
	l := "\n"
	if len(language) > 0 {
		l = language[0] + l
	}

	return "```" + l + text + "\n```"
}

func Code(text string) string {
	return "`" + text + "`"
}

func OrderedList(texts []string) (res string) {
	for i, text := range texts {
		res += strconv.Itoa(i) + ". " + text + "\n"
	}

	res = res[:len(res)-1]

	return
}

func UnorderedList(texts []string) (res string) {
	for _, text := range texts {
		res += "- " + text + "\n"
	}

	res = res[:len(res)-1]

	return
}

func Link(text string, link string) string {
	return "(" + text + ")[" + link + "]"
}

func ImageLink(link string, alt string) string {
	return "!(" + alt + ")[" + link + "]"
}

func Bold(text string) string {
	return "**" + text + "**"
}

func Italic(text string) string {
	return "*" + text + "*"
}

func UnderLine(text string) string {
	return "<u>" + text + "</u>"
}

func DeleteLine(text string) string {
	return "~~" + text + "~~"
}

func Annotation(text string) string {
	return "<!--" + text + "-->"
}

func CuttingLine() string {
	// 分割线也包括"---"和"___"，并且字符之间可以间隔
	return "***"
}

type columnSetting uint

const (
	defaultDisplay columnSetting = iota
	leftJustify
	middleJustify
	rightJustify
)

const (
	structLineString = " | "
	lineFeedString   = "<br />"
	defaultString    = "--"
	leftString       = ":--"
	middleString     = ":--:"
	rightString      = "--:"
)

type FormTable struct {
	Line          uint
	Column        uint
	columnSetting []columnSetting
	data          [][]string
}

func (f *FormTable) Init(line uint, column uint) {
	f.Line = line
	f.Column = column
	f.columnSetting = make([]columnSetting, column)
	f.data = make([][]string, line)
	for i := range f.data {
		f.data[i] = make([]string, column)
	}
}

func (f *FormTable) AddLine(add uint, after ...uint) {
	var (
		af = int(f.Line + add)
		ad = int(add)
	)
	if len(after) > 0 {
		af = int(after[0])
	}

	for i := 0; i < ad; i++ {
		f.data = append(f.data, make([]string, f.Column))
	}

	for i := 0; i < ad; i++ {
		temp := f.data[len(f.data)-1]
		for j := len(f.data) - 1; j > af; j-- {
			f.data[j] = f.data[j-1]
		}
		f.data[af-1] = temp
	}

	f.Line += add
}

func (f *FormTable) AddLineAndSet(texts ...string) {
	f.AddLine(1)

	for i, text := range texts {
		f.SetText(f.Line-1, uint(i), text)
	}
}

func (f *FormTable) AddColumn(add uint, after ...uint) {
	var (
		af = int(add + f.Column)
		ad = int(add)
	)
	if len(after) > 0 {
		af = int(after[0])
	}

	for i, data := range f.data {
		f.data[i] = append(f.data[i], make([]string, ad)...)

		for j := 0; j < ad; j++ {
			temp := data[len(data)-1]
			for k := len(data) - 1; k > af; k-- {
				data[k] = data[k-1]
			}
			data[af-1] = temp
		}
	}

	f.columnSetting = append(f.columnSetting, make([]columnSetting, ad)...)
	for j := 0; j < ad; j++ {
		temp := f.columnSetting[len(f.columnSetting)]
		for k := len(f.columnSetting) - 1; k > af; k-- {
			f.columnSetting[k] = f.columnSetting[k-1]
		}
		f.columnSetting[af-1] = temp
	}

	f.Column += add
}

func (f *FormTable) SetColumnDisplay(column uint, setting columnSetting) {
	f.columnSetting[column] = setting
}

func (f *FormTable) SetText(line uint, column uint, text string) {
	f.data[line][column] = text
}

func (f *FormTable) String() string {
	if f.Line == 0 {
		return ""
	}

	var (
		resB  strings.Builder
		lineB strings.Builder
	)

	// 表格的前一行必须是空格
	resB.WriteByte('\n')
	for i, data := range f.data {
		lineB.Reset()

		writeLine(&lineB, data)
		lineB.WriteByte('\n')

		if i == 0 {
			writeColumnSetting(&lineB, f.columnSetting)
			lineB.WriteByte('\n')
		}

		s := lineB.String()
		resB.WriteString(s)
	}
	// 表格之后一行必须是空格
	resB.WriteByte('\n')

	return resB.String()
}

func writeLine(b *strings.Builder, datas []string) {
	l := len(datas)
	if l == 0 {
		return
	}

	b.WriteString(structLineString[1:])
	for i, data := range datas {
		data = strings.ReplaceAll(data, "\n", lineFeedString)

		b.WriteString(data)

		if i != l-1 {
			b.WriteString(structLineString)
		} else {
			b.WriteString(structLineString[:len(structLineString)-1])
		}
	}
}

func writeColumnSetting(b *strings.Builder, cs []columnSetting) {
	l := len(cs)
	if l == 0 {
		return
	}

	b.WriteString(structLineString[1:])
	for i, setting := range cs {
		switch setting {
		case defaultDisplay:
			b.WriteString(defaultString)

		case leftJustify:
			b.WriteString(leftString)

		case middleJustify:
			b.WriteString(middleString)

		case rightJustify:
			b.WriteString(rightString)

		default:
			b.WriteString(defaultString)

		}

		if i != l-1 {
			b.WriteString(structLineString)
		} else {
			b.WriteString(structLineString[:len(structLineString)-1])
		}
	}
}
