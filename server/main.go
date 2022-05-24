package main

import (
	"fmt"

	"github.com/sunyufeng1/SyfGwfProxy/server/proxy"
)

// linux peizhi  GOOS = linux  GOARCH = amd64
func main() {
	fmt.Println("Socks5 service 启动")
	socks5Obj := new(proxy.ServerObj)
	socks5Obj.Init()
	//
	//var str string
	//fmt.Scan(&str)
	fmt.Println("Socks5 service 进程结束")
}
