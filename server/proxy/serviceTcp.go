package proxy

import (
	"log"
	"net"
)

func (this *ServerObj) initListenerTCP() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listen, err := net.Listen("tcp", ":8092")
	if err != nil {
		log.Panic(err)
	}
	println("tcp 监听启动成功")
	for {
		client, err := listen.Accept()
		if err != nil {
			log.Println(err)

			continue
		}

		go this.handleClientRequestTCP(client)
	}
}

func (this *ServerObj) handleClientRequestTCP(client net.Conn) {

	if client == nil {
		return
	}
	defer client.Close()

	b := make([]byte, 1024)
	bLeng, err := client.Read(b) //第一次从客户端读取请求
	if err != nil {
		log.Println("第一次从客户端读取请求")
		log.Println(err)
		return
	}
	log.Println("tcp 收到连接的请求")
	oldArr := b[:bLeng]
	this.decoder.Decode(b) //步骤1 解密

	if b[0] != 0x05 { //只处理Socks5协议
		//log.Println("不是客户端发起的请求")
		//log.Println(b[0])
		//说明不是通过加密  而是从外部想要进行连接
		this.encoder.Encode(b)
		this.handleServiceRequestTCP(client, oldArr)
		return
	}
	log.Println("tcp 收到代理客户端信息")

	//进入步骤2
	//客户端回应：Socks服务端不需要验证方式
	wiriteBuf := []byte{0x05, 0x00} //第一字节 版本 第二 字节验证方式
	this.encoder.Encode(wiriteBuf)  //步骤1 加密
	client.Write(wiriteBuf)         //第一次回复客户端请求

	log.Println("tcp 向代理客户端发送 认证完成")
	bLen, err := client.Read(b) //第二次读取代理客户端请求 ip和端口
	if err != nil {
		log.Println("第二次读取代理客户端请求 出错")
		//log.Println("阶段2错误 ,%s",err.Error())
		//log.Println("头信息接受 错误")
		return
	}
	log.Println("tcp 接收到 头信息")
	this.decoder.Decode(b) //步骤2 解密

	address := string(b[:bLen])
	log.Println(address)
	//获得了请求的host和port，就开始拨号吧
	//target, err := net.Dial("tcp", address)
	target, err := net.Dial("tcp", address)

	if err != nil {
		log.Println(err.Error() + address)
		log.Println("目标地址连接 失败")
		return
	}

	log.Println("目标地址连接成功")
	println(target.RemoteAddr().String())
	temp := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	this.encoder.Encode(temp)
	client.Write(temp) //响应客户端连接目标服务器成功 第二次回复代理客户端
	log.Println("向客户端表明 目标服务器已经连接成功")
	//client.SetDeadline()
	//进行转发
	go func() {
		err := this.decoder.DecodeCopy(target, client)
		//log.Println("发送到目标地址成功")
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			println("发送到目标的时候超时连接断开1")
			//println(err.Error())
			client.Close()
			target.Close()
		}
	}()

	err = this.encoder.EncodeCopy(client, target) //把发送的目标地址 的内容转发给客户端
	if err != nil {
		//	// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
		println("发送到目标的时候超时连接断开2")
		println(err.Error())
		client.Close()
		target.Close()
	}
	log.Println("阶段线程结束")

}
