package demo

import (
	"context"
	"fmt"
	"github.com/fmyxyz/pod"
	"net"
)

func main() {
	conn, e := net.Dial("tcp", "127.0.0.1:12321")
	if e != nil {
		cc := NewDomeClinetConn(conn)
		cc.Start(conn)
	}
}

func NewDomeClinetConn(conn net.Conn) *pod.ClientConn {
	cc := pod.NewClientConn(0)
	cc.OnReadStart = func(ctx context.Context, conn net.Conn) {
		b := make([]byte, 10, 10)
		conn.Read(b)
		fmt.Println(b)
	}
	cc.OnConnected = func(ctx context.Context, conn net.Conn) {
		fmt.Println("I am in...")
		conn.Write([]byte("I am in..."))
	}
	return cc
}
