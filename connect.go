package pod

import (
	"net"
	"context"
	"sync"
)

type ClientConn struct {
	connid int64
	conn net.Conn
	//同步锁
	mux sync.Mutex
	Lifecycle
}
func (cc *ClientConn)Start(){

}
func NewClientConn(connid int64,conn  net.Conn ) *ClientConn   {
	return &ClientConn{
		connid:connid,
		conn:conn,
	}
}
type ServerConn struct {
	//同步锁
	mux sync.Mutex
	Lifecycle
}

type Lifecycle struct {
	//1、链接建立 :
	OnConnected func(ctx context.Context, conn net.Conn)
	//2、链接关闭 :
	OnConnClose func(ctx context.Context, conn net.Conn) ()
	//3、链接错误：
	OnError func(ctx context.Context, conn net.Conn)
	//7.取消数据传输 :
	OnTransferCancel func(ctx context.Context, conn net.Conn)
	//8.初始化
	init func(ctx context.Context, conn net.Conn)
	//9.销毁
	destroy func(ctx context.Context, conn net.Conn)
	//1.读数据结束
	OnReadEnd func(ctx context.Context, conn net.Conn)
	//2.写数据结束
	OnWriteEnd func(ctx context.Context, conn net.Conn)
	//1.读数据开始
	OnReadStart func(ctx context.Context, conn net.Conn)
	//2.写数据开始
	OnWriteStar func(ctx context.Context, conn net.Conn)
	//1.链接超时 ：
	OnConnTimeout func(ctx context.Context, conn net.Conn)
	//2.响应超时 :
	OnResponseTimeout func(ctx context.Context, conn net.Conn)
	//3.未活动超时 :
	OnActiveTimeout func(ctx context.Context, conn net.Conn)
}
