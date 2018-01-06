package main

import (
	"github.com/fmyxyz/pod"
	"github.com/fmyxyz/pod/utils"
	"log"
	"net"
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

	//sc := pod.NewServerConn(netid)
	//sc.Start(conn)
	var a int
	for {
		a++
		msg := pod.Metadata{}
		err := msg.Deserialize(conn)
		if err != nil {
			log.Println("接收数据错误：", err)
			return
		}

		//注册一系列解码器
		for ; ;  {
			
		}



		log.Println("收到数据：", a, msg)

		/*		sc.PullMessage(&msg)
				log.Println("接收心跳信息：", msg)
				sc.PushMessage(&msg)
				log.Println("发送心跳信息：", msg)*/
	}
}
