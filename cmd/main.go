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
	"flag"
	"fmt"
	"log"
)

func main() {
	AGIDCmd()
}

func AGIDCmd() {
	var (
		dic       string
		file      string
		recursion bool
	)

	flag.StringVar(&dic, "d", "", "需要转换的文件夹，未设置时不生效")
	flag.StringVar(&file, "f", "", "需要转换的文件，未设置时不生效")
	flag.BoolVar(&recursion, "r", false, "是否递归文件夹，默认为false")

	flag.Parse()

	if dic != "" {
		rec := "非递归模式"
		if recursion {
			rec = "递归模式"
		}

		fmt.Printf("正在处理文件夹[%s]:%s\n", rec, dic)

		if err := AGID.LoadAndSaveDic(dic, recursion); err != nil {
			log.Println(err)
			return
		}
	}

	if file != "" {
		fmt.Printf("正在处理文件:%s\n", file)

		if err := AGID.LoadAndSave(file, recursion); err != nil {
			log.Println(err)
			return
		}
	}
}
