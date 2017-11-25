package demo

import (
	"github.com/fmyxyz/pod"
	"net"
	"context"
	"fmt"
)

func main() {
	conn,e:=net.Dial("tcp","127.0.0.1:12321")
	if e!=nil{
		cc:=NewDomeClinetConn(conn)
		cc.Start();
	}
}

func NewDomeClinetConn(conn net.Conn) *pod.ClientConn{
	cc:=pod.NewClientConn(0,conn)
	cc.OnReadStart= func(ctx context.Context, conn net.Conn) {
		b:=make([]byte,10,10)
		conn.Read(b)
		fmt.Println(b)
	}
	cc.OnConnected= func(ctx context.Context, conn net.Conn) {
		fmt.Println("I am in...")
		conn.Write([]byte("I am in..."))
	}
	return cc
}
