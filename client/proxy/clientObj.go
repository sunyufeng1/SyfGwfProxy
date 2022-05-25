package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"github.com/sunyufeng1/SyfGwfProxy/comm"
)

type ClientObj struct {
	encoder *comm.EncoderObj //编码器
	decoder *comm.DecoderObj //解码器

}

type TestObj struct {
	a int
}

func (this *ClientObj) Init() {
	this.encoder = new(comm.EncoderObj)
	this.decoder = new(comm.DecoderObj)
	this.encoder.Init()
	this.decoder.Init()

	this.initListener()
	//this.initListenerTCP()
	//this.initListenerOutTCP()
	//this.initListenerUDP()
}

type ServerIp struct {
	Ip string
}

func (this *ClientObj) initListener() {
	fmt.Println("begin web listen")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listen, err := net.Listen("tcp", ":8099")
	if err != nil {
		log.Panic(err)
	}

	jsonStr := new(ServerIp)

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile("./ServerIP")
	if err != nil {
		print(err.Error())
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, jsonStr)
	if err != nil {
		print(err.Error())
		return
	}

	println(" ServerIp : " + jsonStr.Ip)

	for {
		IE, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("新的连接")
		//_ = IE
		go this.handleIERequest(IE, jsonStr.Ip)
	}

}

func (this *ClientObj) handleIERequest(IE net.Conn, Ip string) {
	if IE == nil {
		return
	}
	defer IE.Close()

	host := "127.0.0.1" //"127.0.0.1"//"47.242.160.106"//Ip //"3.1.24.118"//"3.1.72.147  新加坡 // 日本 3.113.201.225
	port := "8091"
	proxyServer, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Println("连接到代理服务器 %s 失败:%s", net.JoinHostPort(host, port), err)
		return
	}

	defer proxyServer.Close()
	//println("连接代理服务器成功")
	//客户端 ：第一个字节 协议版本 0x05 第二个字节 验证方式 0x00 验证方式占用几个字节 第三个字节 验证方式 0x00 不要验证 0x02 用户名和密码
	wiriteBuf := []byte{0x05, 0x00, 0x00} //说明协议和 验证方式
	this.encoder.Encode(wiriteBuf)
	proxyServer.Write(wiriteBuf)
	//println("向代理服务器发送客户端信息")
	b := make([]byte, 1024)
	num, err := proxyServer.Read(b)
	if err != nil {
		//log.Println(err)
		return
	}
	_ = num

	this.decoder.Decode(b)
	if b[0] != 0x05 || b[1] != 0x00 {
		//log.Println("收到认证信息 错误")
		return
	}
	//log.Println("收到 认证成功信息")

	_, err = IE.Read(b) //读取头信息
	if err != nil {
		//log.Println(err)
		//log.Println("读取头信息 错误")
		return
	}

	headInfo := b[:]
	println(" 头信息*****************")
	//fmt.Printf("%s",string(headInfo[:]))
	log.Printf("str is:%v\n", string(b[:]))
	var method, host1 string
	nIndex := bytes.IndexByte(headInfo[:], '\n')
	if nIndex == -1 {
		return
	}
	fmt.Sscanf(string(headInfo[:nIndex]), "%s%s", &method, &host1)
	//与远程服务端建立成功后 如果是CONNECT的要单独处理
	if method == "CONNECT" { //https 才会走这里 connect 之后才会走post 或者get
		//println("进入了https处理")
		connectByte := []byte("HTTP/1.1 200 Connection established\r\n\r\n")
		IE.Write(connectByte)

		this.encoder.Encode(headInfo)
		proxyServer.Write(headInfo) //这里要把头信息写回去 写给代理服务器
		//log.Println("https 发送头信息给https")
		//fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
		//_ , err = IE.Read(b)
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
		//	//headInfo = b[:]
		//	//log.Println("客户端收到https 第二次的数据")
		//connectByte1 := []byte("HTTP/1.1 200 Connection established\r\n\r\n")
		//IE.Write(connectByte1)
	} else { //这里的话直接走pos 或者get
		this.encoder.Encode(headInfo)
		proxyServer.Write(headInfo) //这里要把头信息写回去 写给代理服务器
	}

	println(" method is :" + method)

	//println("准备服务器连接目标地址的结果")
	_, err = proxyServer.Read(b)
	if err != nil {
		log.Println(err)
		log.Println("无法连接到目标地址 错误1")
		return
	}

	this.decoder.Decode(b)
	//协议正确 却 与真实目标连接成功 才会往下走
	if b[0] != 0x05 || b[1] != 0x00 {
		//log.Println("无法连接到目标地址 错误2")
		return
	}
	println("目标地址连接成功")
	//if method == "CONNECT" {
	//	//this.encoder.Encode(headInfo)
	//	//proxyServer.Write(headInfo)//这里要把头信息写回去 写给代理服务器
	//	//println("https 第二次数据发送完成")
	//}
	go func() {
		err := this.decoder.DecodeCopy(IE, proxyServer)
		if err != nil {
			// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
			println("连接关闭")
			IE.Close()
			//proxyServer.Close()
		}
	}()
	// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	err = this.encoder.EncodeCopy(proxyServer, IE)
	if err != nil {
		// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
		println("连接关闭2")
		IE.Close()
		proxyServer.Close()
	}

}
