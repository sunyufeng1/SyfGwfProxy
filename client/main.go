package main

import (
	"github.com/sunyufeng1/SyfGwfProxy/client/proxy"
)

func main() {
	// 去掉黑色背景 -ldflags="-H windowsgui"

	clientObj := new(proxy.ClientObj)
	clientObj.Init()
}

//情况分析
//1 从客户端发起http https请求
//2 从客户端发起tcp
//3 从客户端发起udp
//4 服务端发起udp
//5 服务端发起tcp

// http https 浏览器设置代理就可以直接发送http请求到指定端口 直接用浏览器测试
// tcp 请求由于太多 又太杂 用代理软件指定某软件的代理窗口 转发到指定的端口
// 测试要开proxifier进行转发 cketTest3 作为发起tcp的软件 Wireshark 查看流量的软件
