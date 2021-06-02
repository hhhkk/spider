package main

import (
	_ "spider/config"
	"spider/log"
	_ "spider/network"
	"spider/cmd"
)

func main() {

	// 释放唯一的文件句柄,平滑退出
	defer log.Close()

	cmd.Start()
}
