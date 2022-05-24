package proxy

import (
	"fmt"
	"log"
	"net"
	"time"

	syf "github.com/sunyufeng1/CommTool/toolElse"
)

func (this *ServerObj) initListenerUDP() {
	fmt.Println("begin udp listen")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	service := ":8095"

	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	syf.CheckError(err)

	clientConn, err := net.ListenUDP("udp", udpAddr)
	defer clientConn.Close()
	if syf.CheckErrorResult(err) {
		return
	}

	for {
		this.handleUDPRequest(clientConn)
	}
}
func (this *ServerObj) handleUDPRequest(clientConn *net.UDPConn) {
	if clientConn == nil {
		return
	}

	//println("向代理服务器发送客户端信息")
	b := make([]byte, 1024)
	num, addr, err := clientConn.ReadFromUDP(b) //第一次从代理服务器中读取

	if syf.CheckErrorResult(err, "第一次从代理服务器中读取 :") {
		return
	}
	_ = num

	//this.decoder.Decode(b);
	//if b[0] != 0x05 || b[1] != 0x00 { //第一次从代理服务器收到的解析内容  自己设置的
	//	//log.Println("收到认证信息 错误")
	//	return
	//}
	println(string(b[:num]))

	//wiriteBuf := []byte{0x05, 0x00} //第一字节 版本 第二 字节验证方式
	//this.encoder.Encode(wiriteBuf)//步骤1 加密
	clientConn.WriteToUDP([]byte(time.Now().String()), addr) //第一次回复客户端请求
	log.Println("发送 认证成功信息 ")
	//开始处理客服端收到的认证请求

	////客户端 ：第一个字节 协议版本 0x05 第二个字节 验证方式 0x00 验证方式占用几个字节 第三个字节 验证方式 0x00 不要验证 0x02 用户名和密码
	//_, err = clientConn.Read(b) //读取头信息 第一次客户端数据开始
	//if  err != nil {
	//	log.Print("第一次客户端数据开始 错误 ")
	//	log.Println(err)
	//	return
	//}
	//
	////headInfo := b[:]
	//println(" 头信息*****************")
	////println(headInfo)
	////fmt.Printf("%s",string(headInfo[:]))
	//log.Printf("str is:%v\n", b[:])
	//tcpp.Write([]byte{0x05, 0x00}) //第一次返回客户端信息 认证
	//log.Println("跟客户端的认证完成")
	////客户端收到的认证请求处理结束
	//
	////客户端收到的连接请求开始处理
	//lenL, err := tcpp.Read(b) //读取socks5具体数据 第二次读取客户端
	//if err != nil {
	//	log.Print("sock5具体数据  : ")
	//	log.Println(err)
	//	return
	//}
	////目标地址开始解析
	//log.Printf("socks 5 detail is:%v\n", b[:])
	//addrType := b[3] //0x01 IPv4地址，DST.ADDR部分4字节长度  0x04 IPv6地址，16个字节长度 0x03域名，DST ADDR部分第一个字节为域名长度，DST.ADDR剩余的内容为域名，没有\0结尾。
	//log.Print("IP地址类型 :")
	//log.Println(addrType)
	//addB := []byte{}
	////portB := []byte{}
	//portB := b[lenL-2 : lenL]
	//targetPort := binary.BigEndian.Uint16(portB) //大端序模式的字节转为int32
	//log.Printf("bytes to int32: %d\n\n", targetPort)
	//var targetIp net.IP
	//if (addrType == 0x01) { //0x01 IPv4地址，DST.ADDR部分4字节长度
	//	addB = b[4 : lenL-2]
	//	targetIp = net.IPv4(addB[0], addB[1], addB[2], addB[3])
	//} else if (addrType == 0x04) { // 0x04 IPv6地址，16个字节长度
	//	addB := b[4 : lenL-2]
	//	targetIp = net.ParseIP(string(addB))
	//} else if (addrType == 0x03) { //0x03域名，DST ADDR部分第一个字节为域名长度，DST.ADDR剩余的内容为域名，没有\0结尾。
	//	//domainLen := b[4]
	//	addB = b[5 : lenL-2]
	//	if addr, er := net.ResolveIPAddr("ip", string(addB)); er != nil {
	//		log.Println("域名解析失败")
	//		return
	//	} else {
	//		targetIp = addr.IP
	//	}
	//} else {
	//	println("真实目标地址错误")
	//	return
	//}
	//log.Printf("net.ParseIP bytes to int32: %v\n\n", targetIp)
	////目标地址解析结束
	//
	////发送地址和IP给代理服务器 准备第二次连接代理服务器
	//ipAndPort := targetIp.String() + ":" + strconv.Itoa(int(targetPort))
	//log.Print("目标地址是   ")
	//log.Println(ipAndPort)
	//ipAndPortB := []byte(ipAndPort)
	//this.encoder.Encode(ipAndPortB)
	//proxyServer.Write(ipAndPortB) //ip地址和端口 写给代理服务器
	////println("准备服务器连接目标地址的结果")
	//_, err = proxyServer.Read(b) //第二次从代理服务器得到回复
	//if err != nil {
	//	log.Println(err)
	//	log.Println("无法连接到目标地址 错误1")
	//	return
	//}
	//
	//this.decoder.Decode(b)
	////协议正确 却 与真实目标连接成功 才会往下走
	//if b[0] != 0x05 || b[1] != 0x00 {
	//	//log.Println("无法连接到目标地址 错误2")
	//	return
	//}
	//println("目标地址连接成功")
	//
	////发送ip和地址给代理服务器结束
	//
	////将连接结果告诉目标客户端
	//tcpp.Write([]byte{0x05, 0x00, 0x01, 0x01, 0x7f, 0x00, 0x00, 0x01, 31, 156}) //第二次返回客户端信息 连接   对外 127.0.0.1：8087
	//log.Print(" connect finish")
	////log.Print(binary.BigEndian.Uint16([]byte{31,156}))
	////0x05 SOCKS5协议版本
	////0x00 连接成功
	////0x01 RSV保留字段
	////0x01 地址类型为IPV4
	////0x7f 0x00 0x00 0x01 代理服务器连接目标服务器成功后的代理服务器IP, 127.0.0.1
	////0xaa 0xaa 代理服务器连接目标服务器成功后的代理服务器端口（代理服务器使用该端口与目标服务器通信），本例端口号为43690
	//
	////准备进行转发
	//
	//go func() {
	//	err := this.decoder.DecodeCopy(tcpp, proxyServer)
	//	if err != nil {
	//		// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
	//		println("转发 svnc 连接关闭")
	//		tcpp.Close()
	//		proxyServer.Close()
	//	}
	//}()
	//// 从 localUser 发送数据发送到 proxyServer，这里因为处在翻墙阶段出现网络错误的概率更大
	//err = this.encoder.EncodeCopy(proxyServer, tcpp)
	//if err != nil {
	//	// 在 copy 的过程中可能会存在网络超时等 error 被 return，只要有一个发生了错误就退出本次工作
	//	println("转发 svnc 连接关闭2")
	//	tcpp.Close()
	//	proxyServer.Close()
	//}

}
