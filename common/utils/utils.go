package utils

import (
	"awesomeProject/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func ReadPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("等待读取客户端发送的数据...")
	//conn.Read 在conn没有被关闭的情况下，才会阻塞
	//如果客户端关闭里conn 则，就不会阻塞
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("conn.Read error")
		return
	}
	//根据bug[:4] 转成一个unit32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4])

	//根据pkgLen 读取消息内容
	var n int
	n, err = conn.Read(buf[:pkgLen]) //从conn里读取pkgLen个字节扔到buf这个切片(缓存)里
	if uint32(n) != pkgLen || err != nil {
		//fmt.Println("conn.Read fail err=", err)
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

func WritePkg(conn net.Conn, data []byte) (err error) {
	//先发送一个长度给对方
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
	//发送data本身
	n, err = conn.Write(data)
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("conn.write(byte) fail", err)
		return
	}
	return err
}
