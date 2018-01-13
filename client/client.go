package main

import (
	"github.com/fmyxyz/pod"
	//"github.com/fmyxyz/pod/utils"
	//"log"
	"net"
	"fmt"
	"log"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	defer conn.Close()
	//netid := utils.GenerationId()
	//var cc = pod.NewClientConn(netid)
	//cc.Start(conn)
	msg := pod.NewHeatbeatMsg()

	bmsg := pod.NewBinaryMsg()

	textmsg := pod.NewTextMsg()
	for i := 0; i < 10; i++ {
		err := msg.Serialize(conn)
		if err != nil {
			return
		}

		bmsg.Data = []byte(fmt.Sprintln("BinaryMsg:", i))
		bmsg.Serialize(conn)
		log.Println("发送完成消息：",bmsg)

		textmsg.Data=fmt.Sprintln("text 数据：", i)
		textmsg.Serialize(conn)
		//cc.PushMessage(&msg)
		//log.Println("发送消息：", msg)
		//cc.PullMessage(&msg)
		//log.Println("接收消息：", msg)
		//<-time.After(msg.Duration)
	}
}
