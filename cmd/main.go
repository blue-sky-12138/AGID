/*
@Description:
@File : main
@Author : blue_sky_12138
@Version: 1.0.0
@Date : 2021/11/29 21:54
*/

package main

import (
	"AGID"
	"fmt"
)

func main() {
	err := AGID.LoadAndSave(`C:\Users\Lenovo\Desktop\脚本测试集`, false)
	fmt.Println(err)
}
