package comm

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

type PasswordObj struct {
	content  []byte
	bookName string
}

func (this *PasswordObj) Init() {
	//本地是否有查找到密码本
	//如果没有密码本则生成密码本
	//如果有则进行本地文件的读取
	this.bookName = "./book.txt"
	bookData := this.readBook()
	if len(bookData) == 0 {
		bookData = this.createBook()
		//fmt.Println("not find ")
	}
	//fmt.Println("ok ")
	this.content = bookData
}

//生成密码本
func (this *PasswordObj) createBook() []byte {
	bookData := this.beginRandom()
	this.saveBook(bookData)
	return bookData
}

//随机256个出来
func (this *PasswordObj) beginRandom() []byte {
	intArr := rand.Perm(256)
	byteArr := make([]byte, 256)
	for i := 0; i <= 255; i++ {
		if i == intArr[i] {
			return this.beginRandom()
		}
		byteArr[i] = byte(intArr[i])
		fmt.Println(byte(intArr[i]))
	}
	return byteArr
}

//将密码本保存到本地文件
func (this *PasswordObj) saveBook(data []byte) {
	err2 := ioutil.WriteFile(this.bookName, data, 0666) //写入文件(字节数组)
	this.check(err2)
}

//读取本地密码本
func (this *PasswordObj) readBook() []byte {
	result := make([]byte, 0)
	find := this.checkFileIsExist(this.bookName)
	if find == true {
		data, err := ioutil.ReadFile(this.bookName)
		this.check(err)
		result = data
	}

	return result
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func (this *PasswordObj) checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func (this *PasswordObj) check(e error) {
	if e != nil {
		panic(e)
	}
}

func (this *PasswordObj) GetContent() []byte {
	return this.content
}
