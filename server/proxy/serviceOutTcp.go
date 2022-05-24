package proxy

import (
	"log"
	"net"
)

//从外部来的连接
func (this *ServerObj) handleServiceRequestTCP(targetService net.Conn, bArr []byte) {

	if targetService == nil {
		return
	}
	defer targetService.Close()
	remoteAddr := targetService.LocalAddr().String()
	targetService.RemoteAddr()
	// 还要分离出地址和端口
	this.handleServerTCPRequest(targetService, "127.0.0.1", "8087", remoteAddr, bArr)
}

func (this *ServerObj) handleServerTCPRequest(targetService net.Conn, host string, port string, remoteAddr string, bArr []byte) {
	if targetService == nil {
		return
	}
	defer targetService.Close()
	//开始链接客户端
	proxyClient, err := net.Dial("tcp", net.JoinHostPort(host, "8096"))
	if err != nil {
		log.Println("tcp 连接到代理客户端 %s 失败:%s", net.JoinHostPort(host, "8096"), err)
		return
	}

	defer proxyClient.Close()
	log.Println("连接代理客户端成功")

	wiriteBuf := []byte{0x05, 0x00, 0x00} //说明协议和 验证方式
	this.encoder.Encode(wiriteBuf)
	proxyClient.Write(wiriteBuf) //第一次写给代理服务器
	b := make([]byte, 1024)
	num, err := proxyClient.Read(b) //第一次从代理服务器中读取
	if err != nil {
		log.Print("第一次从代理服务器中读取 :")
		log.Println(err)
		return
	}
	_ = num
	this.decoder.Decode(b)
	if b[0] != 0x05 || b[1] != 0x00 { //第一次从代理服务器收到的解析内容  自己设置的
		log.Println("收到认证信息 错误")
		return
	}
	log.Println("收到 认证成功信息 ")
	//开始处理客服端收到的认证请求

	//进行ip的发送
	println(" 目标客户端地址" + remoteAddr)
	wiriteBuf = []byte(remoteAddr)
	this.encoder.Encode(wiriteBuf)
	proxyClient.Write(wiriteBuf)

	//等待连接目标客户端连接成功
	num, err = proxyClient.Read(b) //
	if err != nil {
		log.Print("等待连接目标客户端连接失败")
		log.Println(err)
		return
	}
	_ = num
	this.decoder.Decode(b)
	if b[0] != 0x05 || b[1] != 0x00 { //
		log.Println("等待连接目标客户端连接失败 2")
		return
	}
	log.Println("连接目标客户端连接成功 ")
	wiriteBuf = bArr
	this.encoder.Encode(wiriteBuf)
	proxyClient.Write(wiriteBuf)

	//开始接收转发
	go func() {
		err := this.decoder.DecodeCopy(targetService, proxyClient)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			println("转发 svnc 连接关闭")
			targetService.Close()
			proxyClient.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	err = this.encoder.EncodeCopy(proxyClient, targetService)
	if err != nil {
		// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
		println("转发 svnc 连接关闭2")
		targetService.Close()
		proxyClient.Close()
	}

}
