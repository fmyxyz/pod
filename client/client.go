package main

import (
	"github.com/fmyxyz/pod"
	"github.com/fmyxyz/pod/utils"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	defer conn.Close()
	netid := utils.GenerationId()
	var cc = pod.NewClientConn(netid)
	cc.Start(conn)
	msg := pod.NewHeatbeatMsg()
	msg.MsgType = 1
	msg.Duration = 3 * time.Second
	for i := 0; i < 100; i++ {
		msg.Serialize(conn)
		log.Println("发送完成消息：", msg)
		//cc.PushMessage(&msg)
		//log.Println("发送消息：", msg)
		//cc.PullMessage(&msg)
		//log.Println("接收消息：", msg)
		//<-time.After(msg.Duration)
	}
}
