package main

import (
	//"github.com/fmyxyz/pod"
	"github.com/fmyxyz/pod/utils"
	"log"
	"net"
	"github.com/fmyxyz/pod"
)

func main() {
	l, e := net.Listen("tcp", "127.0.0.1:8080")
	if e != nil {
		return
	}
	defer l.Close()
	log.Println("本地侦听:", l.Addr())
	for {
		c, er := l.Accept()
		if er != nil {
			break
		}
		go handleConn(c)
	}
}
func handleConn(conn net.Conn) {
	log.Println("新的连接：", conn.RemoteAddr())
	netid := utils.GenerationId()
	log.Println("netid:", netid)
	//心跳信息
	hitmsg := pod.NewHeatbeatMsg()
	//注册
	regMsg := pod.RegistryMessage(pod.HEATBEAT, &hitmsg)

	var a int
	for {
		a++
		//b:=make([]byte,17)
		//i,err:=conn.Read(b)
		//log.Println(b,i)
		msg := pod.Metadata{}
		err := msg.Deserialize(conn)
		if err != nil {
			log.Println("接收数据错误：", err)
			return
		}
		//设置获取的元数据
		regMsg.SetMetadata(msg)
		err = regMsg.Deserialize(conn)

		if err != nil {
			log.Println("接收数据错误：", err)
			return
		}

		//sc.PullMessage(&msg)
		//log.Println("接收心跳信息：", msg)
		//sc.PushMessage(&msg)
		//log.Println("发送心跳信息：", msg)
	}
}
