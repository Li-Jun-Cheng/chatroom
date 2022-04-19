package main

import (
	"awesomeProject/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//写一个函数，完成登录

func login(userId int, userPwd string) (err error) {
	//下一个就要开始定协议...
	fmt.Printf("userId=%d userPwd=%s", userId, userPwd)

	//1.连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889") //拨号
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()
	//2.准备通过conn发消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4.将loginMes 序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//5.把data赋给mes.Data字段
	mes.Data = string(data)

	//6.将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("", err)
		return
	}
	//7这个时候data就是我们要发送的消息
	//7.1先把data的长度发送给服务器
	//先获取到data的长度->转成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:])
	if n != 4 || err != nil {
		fmt.Println("conn.write(byte) fail", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	//这里还需要处理服务器返回的消息。
	return
}
