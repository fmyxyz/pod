package pod

type MsgaAction interface {
	//发送消息
	PushMessage(message Message)
	//获取消息
	PullMessage() Message
	//取消
	Cancel()
	//关闭
	Close()
}
