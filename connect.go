package pod

import (
	"net"
	"context"
	"sync"
	"bufio"
	"time"
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
	buReader *bufio.Reader
	buWriter *bufio.Writer
	//同步锁
	mux sync.Mutex
	Lifecycle
	//1.链接超时 ：
	connTimeout time.Duration
	//2.响应超时 :
	responseTimeout time.Duration
	//3.未活动超时 :
	activeTimeout time.Duration
}
//设置链接超时时间
func (bc *BaseConn)SetConnTimeout(connTimeout time.Duration){
	bc.connTimeout=connTimeout
}
//设置相应超时时间
func (bc *BaseConn)SetResponseTimeout(responseTimeout time.Duration){
	bc.responseTimeout=responseTimeout
}
//设置活动超时时间
func (bc *BaseConn)SetActiveTimeout(activeTimeout time.Duration){
	bc.activeTimeout=activeTimeout
}
func (bc *BaseConn)GetTimeout() time.Duration{
	return  bc.activeTimeout
}
//启动
func (bc *BaseConn) Start(conn net.Conn){
	bc.buReader=bufio.NewReader(conn)
	bc.buWriter=bufio.NewWriter(conn)
}
//发送消息
func (bc *BaseConn) PushMessage(message *Message){
	bc.mux.Lock()
	defer bc.mux.Unlock()
	bc.OnWriteStar(context.TODO(),bc.conn)
	//写入数据
	bc.buWriter.Write((*message).Serialize())
	bc.OnWriteEnd(context.TODO(),bc.conn)
}
//获取消息
//len 字节数
//message 消息实例
func (bc *BaseConn)PullMessage(len func()int64,message *Message) *Message{
	bc.mux.Lock()
	defer bc.mux.Unlock()
	bc.OnReadStart(context.TODO(),bc.conn)
	b:=make([]byte,len())
	bc.buReader.Read(b)
	(*message).Deserialize(b)
	bc.OnReadEnd(context.TODO(),bc.conn)
	return message
}
//取消
func (bc *BaseConn)Cancel(){

}
//关闭
func (bc *BaseConn)Close(){

}

