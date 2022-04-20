package main

import (
	"awesomeProject/common/message"
	"awesomeProject/common/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

//编写一个函数serverProcessLogin函数，专门处理登录请求
func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
	//核心代码
	//1.先从mes中取出mes.Data,并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//声明一个resMes,用于向客户端返回状态信息
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//2再声明一个LoginResMes，并完成赋值
	var loginResMes message.LoginResMes

	//如果用户的id为100  密码=123456，认为合法，否则不合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		//合法
		loginResMes.Code = 200
		loginResMes.Message = "非常棒！你知道吗？"

	} else {
		//不合法
		loginResMes.Code = 500 //500状态码表示该用户不存在
		loginResMes.Error = "该用户不存在，请您注册！！！"
	}
	//3将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.marshal fail", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.对resMes进行序列化，准备发送
	data, err = json.Marshal(resMes) //整体序列化
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	//6.发送data,我们将其封装到writePkg()函数中
	err = utils.WritePkg(conn, data)
	return err
}

//编写一个ServerProcessMes函数
//功能：根据客户端发送的消息种类不同，决定调用哪个函数来处理
func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登录的逻辑
		err = serverProcessLogin(conn, mes)
	case message.RegisterMesType:
	//处理注册
	default:
		fmt.Println("消息类型不存在，无法处理...")
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
		mes, err := utils.ReadPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器端也退出...")
				return
			} else {
				err = errors.New("process error")
				return
			}
		}
		fmt.Println("mes=", mes)
		serverProcessMes(conn, &mes)
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
