package pod
//消息
type Message interface {
	//序列化
	Serialize() []byte
	//反序列化
	Deserialize([]byte)
}

type StringMsg string

func (sm *StringMsg) Serialize() []byte {
	return []byte(*sm)
}

func (sm *StringMsg) Deserialize(bt []byte) {
	*sm = StringMsg(bt)
}