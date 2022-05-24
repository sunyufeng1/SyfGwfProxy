package proxy

import (
	"fmt"
	"log"
	"net"
)

func (this *ClientObj) initListenerOutTCP() {
	fmt.Println("begin out tcp listen")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listen, err := net.Listen("tcp", ":8096")
	if err != nil {
		log.Panic(err)
	}

	for {
		tcpp, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("新的连接")
		//_ = IE

		go this.handleOutTCPRequest(tcpp)
	}
}

func (this *ClientObj) handleOutTCPRequest(tcpp net.Conn) {
	if tcpp == nil {
		return
	}
	defer tcpp.Close()

	b := make([]byte, 1024)
	_, err := tcpp.Read(b) //第一次从代理服务器读取请求
	if err != nil {
		log.Println("第一次从代理服务器读取请求")
		log.Println(err)
		return
	}
	log.Println("tcp 收到连接的请求")
	this.decoder.Decode(b) //步骤1 解密
	if b[0] != 0x05 {      //只处理Socks5协议
		//log.Println("不是客户端发起的请求")
		//log.Println(b[0])
		//说明不是通过加密  而是从外部想要进行连接
		return
	}
	log.Println("tcp 收到代理服务端信息")

	wiriteBuf := []byte{0x05, 0x00} //第一字节 版本 第二 字节验证方式
	this.encoder.Encode(wiriteBuf)  //步骤1 加密
	tcpp.Write(wiriteBuf)           //第一次回复代理服务器请求

	log.Println("tcp 向代理服务器发送 认证完成")

	//接收ip 并且开始连接
	//客户端收到的连接请求开始处理
	len, rerr := tcpp.Read(b) //读取socks5具体数据 第二次读取客户端
	if rerr != nil {
		log.Print("sock5具体数据  : ")
		log.Println(rerr)
		return
	}
	this.decoder.Decode(b)
	//目标地址开始解析
	log.Printf("socks 5 detail is:%v\n", b[:len])
	ipString := string(b[:len])
	targetPort := "8087" //strings.Split( ipString ,":")[1]
	log.Println(targetPort)
	log.Println(ipString)

	//开始连接目标客户端
	targetClient, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", targetPort))
	if err != nil {
		log.Println("tcp 连接到代理客户端 %s 失败:%s", net.JoinHostPort("127.0.0.1", targetPort), err)
		return
	}
	defer targetClient.Close()
	wiriteBuf = []byte{0x05, 0x00}
	this.encoder.Encode(wiriteBuf)
	tcpp.Write(wiriteBuf) //通知连接目标客户端成功
	//发送ip和地址给代理服务器结束

	//开始接收转发
	go func() {
		err := this.decoder.DecodeCopy(targetClient, tcpp)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			println("转发 svnc 连接关闭")
			tcpp.Close()
			targetClient.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	err = this.encoder.EncodeCopy(tcpp, targetClient)
	if err != nil {
		// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
		println("转发 svnc 连接关闭2")
		tcpp.Close()
		targetClient.Close()
	}
}
