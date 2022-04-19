package main

import (
	"awesomeProject/common/message"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("等待读取客户端发送的数据...")
	_, err = conn.Read(buf[:4])
	if err != nil {
		err = errors.New("conn.Read error")
		return
	}
	//根据bug[:4] 转成一个unit32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])

	//根据pkgLen 读取消息内容
	var n int
	n, err = conn.Read(buf[:pkgLen]) //从conn里读取pkgLen个字节扔到buf这个切片(缓存)里
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.Read fail err=", err)
		return
	}
	//把pkg反序列化成 -> message.Message
	json.Unmarshal(buf[:pkgLen], &mes) //这里必须传mes的地址使用&,如果没有不传地址，拿到的mes就是空的
	if err != nil {
		fmt.Println("json.Unmarsha err=", err)
		return
	}

	return
}

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里需要延时关闭conn
	defer conn.Close()
	//循环的读取客户端发送的信息
	for {
		//这里我们将读取数据包，直接封装成一个函数readPkg(),返回Message,Err
		mes, err := readPkg(conn)
		if err != nil {
			err = errors.New("process error")
			return
		}
		fmt.Println("mes=", mes)
	}
}

func main() {
	//提示信息
	fmt.Println("服务器在8889端口监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	//一旦监听成功，就等待客户端来连接服务器
	for {
		fmt.Println("等待客户端来连接服务器....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//一旦连接成功，则启动一个协程和客户端保持通讯...
		go process(conn)
	}
}
