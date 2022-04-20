package main

import (
	"awesomeProject/common/message"
	"awesomeProject/common/utils"
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
	var mes message.Message         //发送消息的结构体、封装数据的依托
	mes.Type = message.LoginMesType //这里说明我们发送的是什么类型的消息
	//3.创建一个LoginMes结构体
	var loginMes message.LoginMes //定义结构体相当于创建一个对象
	loginMes.UserId = userId      //给结构体的属性赋值
	loginMes.UserPwd = userPwd

	//4.将loginMes 序列化
	data, err := json.Marshal(loginMes) //值封装好以后要序列化才能通过网络进行传输
	if err != nil {
		fmt.Println("json.Marshal err=", err) //序列化出错了
		return
	}
	//5.把data赋给mes.Data字段
	mes.Data = string(data) //把序列化后的值赋值的消息结构体的Data属性

	//6.将mes进行序列化
	data, err = json.Marshal(mes) //再将消息结构体本身也序列化，序列化一定要彻底，层层序列化
	if err != nil {
		fmt.Println("", err) //序列化出错了
		return
	}
	//7这个时候data就是我们要发送的消息
	//7.1先把data的长度发送给服务器
	//先获取到data的长度->转成一个表示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data)) //先得到要发送消息结构体的长度
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], pkgLen) //将长度信息转换为byte类型的切片，因为下面的conn.Write()函数的入参的数据类型是byte类型的切片
	//发送长度
	n, err := conn.Write(buf[:]) //写入就是发送到服务器那边了，返回值中n就返回的长度信息了
	if n != 4 || err != nil {
		fmt.Println("conn.write(byte) fail", err)
		return
	}
	//代码运行到这里说明网络联通性不成问题
	_, err = conn.Write(data) //正式发送消息
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	//time.Sleep(2 * time.Second)
	//fmt.Println("休眠了2s...")
	//消息发送成功，这里还需要处理服务器返回的消息。
	mes, err = utils.ReadPkg(conn) //mes就是
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
	}
	//将mes的data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//fmt.Println("登录成功")
		fmt.Printf("返回状态码：%d，返回值：%s", loginResMes.Code, loginResMes.Message)
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}
	return
}
