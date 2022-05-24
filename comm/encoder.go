package comm

import (
	"io"
	"log"
	"net"
)

type EncoderObj struct {
	bookData []byte
}

func (this *EncoderObj) Init() {
	bookObj := new(PasswordObj)
	bookObj.Init()
	this.bookData = bookObj.GetContent()
}

func (this *EncoderObj) Encode(oriData []byte) {
	for key, value := range oriData {
		//fmt.Println(key)
		oriData[key] = this.bookData[value]
	}
}

func (this *EncoderObj) EncodeSend(conn net.Conn, oriData []byte) (n int, err error) {
	n, err = conn.Read(oriData)
	if err != nil {
		return n, err
	}
	this.Encode(oriData)
	return n, err
}

func (this *EncoderObj) EncodeCopy(target net.Conn, src net.Conn) error {
	buf := make([]byte, 1024)
	for {
		readCount, errRead := src.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			log.Println("有数据被加密转发")
			writeBuf := buf[0:readCount]
			log.Printf("%s", string(writeBuf[:]))
			this.Encode(writeBuf)
			writeCount, errWrite := target.Write(writeBuf)
			log.Println("加密转发 ", writeCount)
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				log.Println("加密前后 数据长度不一致")
				return io.ErrShortWrite
			}
		}
	}
}
