package proxy

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/sunyufeng1/SyfGwfProxy/comm"
)

type ServerObj struct {
	encoder *comm.EncoderObj //编码器
	decoder *comm.DecoderObj //解码器
}

func (this *ServerObj) Init() {
	this.encoder = new(comm.EncoderObj)
	this.decoder = new(comm.DecoderObj)
	this.encoder.Init()
	this.decoder.Init()

	//this.initListener()
	// this.initListenerTCP()
	this.initListenerUDP()
}

func (this *ServerObj) initListener() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listen, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Panic(err)
	}

	for {
		client, err := listen.Accept()
		if err != nil {
			log.Println(err)

			continue
		}

		go this.handleClientRequest(client)
	}
}

//有加密
func (this *ServerObj) handleClientRequest(client net.Conn) {

	if client == nil {
		return
	}
	defer client.Close()

	b := make([]byte, 1024)
	_, err := client.Read(b)
	if err != nil {
		//log.Println(err)
		return
	}
	log.Println("收到连接的请求")
	this.decoder.Decode(b) //步骤1 解密
	if b[0] != 0x05 {      //只处理Socks5协议
		//log.Println("不是客户端发起的请求")
		//log.Println(b[0])
		return
	}
	log.Println("收到客户端信息")
	//进入步骤2
	//客户端回应：Socks服务端不需要验证方式
	wiriteBuf := []byte{0x05, 0x00} //第一字节 版本 第二 字节验证方式
	this.encoder.Encode(wiriteBuf)  //步骤1 加密
	client.Write(wiriteBuf)

	log.Println("向客户端发送 认证完成")
	_, err = client.Read(b)
	if err != nil {
		//log.Println("阶段2错误 ,%s",err.Error())
		//log.Println("头信息接受 错误")
		return
	}
	log.Println("接收到 头信息")
	this.decoder.Decode(b) //步骤2 解密

	var method, host, address string
	nIndex := bytes.IndexByte(b[:], '\n')
	if nIndex == -1 {
		return
	}
	fmt.Sscanf(string(b[:nIndex]), "%s%s", &method, &host)
	log.Println(method)
	log.Println(host)
	hostPortURL, err := url.Parse(host)
	if err != nil {
		log.Println(err.Error() + host)
		return
	}
	log.Println("准备拼接地址")
	if hostPortURL.Scheme != "http" { //https访问 443
		address = hostPortURL.Scheme + ":" + hostPortURL.Opaque
		log.Println(1)
		log.Println(hostPortURL.Opaque)
	} else { //http访问
		if strings.Index(hostPortURL.Host, ":") == -1 { //host不带端口， 默认80
			address = hostPortURL.Host + ":80"
			log.Println(2)
		} else {
			address = hostPortURL.Host
			log.Println(3)
		}
	}
	log.Println(address)
	//获得了请求的host和port，就开始拨号吧
	target, err := net.Dial("tcp", address)
	if err != nil {
		//log.Println(err.Error() + address)
		//log.Println("目标地址连接 失败")
		return
	}
	log.Println("目标地址连接成功")
	temp := []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	this.encoder.Encode(temp)
	client.Write(temp) //响应客户端连接目标服务器成功
	//log.Println("向客户端表明 目标服务器已经连接成功")
	//与远程服务端建立成功后 如果是CONNECT的要单独处理
	if method == "CONNECT" { //https 才会走这里 connect 之后才会走post 或者get
		//客户端已经处理
		//connectByte := []byte("HTTP/1.1 200 Connection established\r\n\r\n")
		//this.encoder.Encode(connectByte)
		//client.Write(connectByte)
		/////fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else { //这里的话直接走pos 或者get
		target.Write(b) //这里要把头信息写回去 写给目标地址
	}

	////进行转发
	go func() {
		err := this.decoder.DecodeCopy(target, client)
		//log.Println("发送到目标地址成功")
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			//println("发送到目标的时候超时连接断开")
			//println(err.Error())
			client.Close()
			target.Close()
		}
	}()

	err = this.encoder.EncodeCopy(client, target) //把发送的目标地址 的内容转发给客户端
	if err != nil {
		// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
		//println("发送到目标的时候超时连接断开2")
		//println(err.Error())
		client.Close()
		target.Close()
	}
	//log.Println("阶段线程结束")
}
