/*
@Description:
@File : util
@Author : blue_sky_12138
@Version: 1.0.0
@Date : 2021/12/18 15:12
*/

package AGID

import (
	"AGID/markdown"
	"reflect"
	"strings"
)

func writeWithNewLine(b *strings.Builder, text string) {
	b.WriteString(text)
	b.WriteString("\n\n")
}

func writeWithEnter(b *strings.Builder, text string) {
	b.WriteString(text)
	b.WriteByte('\n')
}

func newDefaultFormTable() (ft markdown.FormTable) {
	ft.Init(1, 4)

	ft.SetText(0, 0, "字段名")
	ft.SetText(0, 1, "必选")
	ft.SetText(0, 2, "类型")
	ft.SetText(0, 3, "说明")

	return
}

func isZero(i interface{}) bool {
	return reflect.ValueOf(i).IsZero()
}
