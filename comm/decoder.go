package comm

import (
	"io"
	"log"
	"net"
)

type DecoderObj struct {
	bookData []byte
}

func (this *DecoderObj) Init() {
	bookObj := new(PasswordObj)
	bookObj.Init()
	this.bookData = bookObj.GetContent()
}

func (this *DecoderObj) Decode(oriData []byte) {
	for key, value := range oriData {
		for bookKey, bookValue := range this.bookData {
			if bookValue == value {
				oriData[key] = byte(bookKey)
				break
			}
		}
	}
}

func (this *DecoderObj) DecodeRead(conn net.Conn, buf []byte) (n int, err error) {
	n, err = conn.Read(buf)
	if err != nil {
		return n, err
	}
	this.Decode(buf)
	return n, err
}

func (this *DecoderObj) DecodeCopy(target net.Conn, src net.Conn) error {
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
			log.Println("decode写入的长度是有")
			writeBuf := buf[0:readCount]
			this.Decode(writeBuf)
			writeCount, errWrite := target.Write(writeBuf)

			log.Println("decode写入的长度是", writeCount)
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				log.Println("DecodeCopy 加密解密前后长度不一致")
				return io.ErrShortWrite
			}
			log.Printf("%s", string(writeBuf[:]))
		}
	}
}
