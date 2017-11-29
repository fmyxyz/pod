package pod

type ConnAction interface {
	//发送消息
	PushMessage(message Message)
	//获取消息
	PullMessage(len func()int64,message Message) Message
	//取消
	Cancel()
	//关闭
	Close()
}
