package pod

import (
	"net"
	"context"
	"sync"
)

//生命周期
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

//基本连接
type BaseConn struct {
	connid int64
	conn net.Conn
	//同步锁
	mux sync.Mutex
	Lifecycle
}

//发送消息
func (bc *BaseConn) PushMessage(message Message){
	bc.conn.Write(message.Serialize())
}
//获取消息
func (bc *BaseConn)PullMessage(len func()int64,message Message) Message{
	l:=len()
	b:=make([]byte,l)
	bc.conn.Read(b)
	message.Deserialize(b)
	return message
}
//取消
func (bc *BaseConn)Cancel(){

}
//关闭
func (bc *BaseConn)Close(){

}

