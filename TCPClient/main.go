package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	//showTCPClient()
	showClientStickyBuns()
}

func EncodeSB(msg string) ([]byte, error) {
	//读取消息,转为int32
	length := int32(len(msg))
	pkg := new(bytes.Buffer)
	//写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	//写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(msg))
	if err != nil {
		return nil, err
	}
	//返回
	return pkg.Bytes(), nil
}

func showTCPClient() {
	//1.与server建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	bufoRear := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("请说话:")
		msg, _ := bufoRear.ReadString('\n')
		if msg == "exit" {
			break
		}
		//2,通信 发送数据
		conn.Write([]byte(msg))
	}
	conn.Close()
}

func showClientStickyBuns() {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		data, err := EncodeSB(msg)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
		conn.Write(data)
	}
}
