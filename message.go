package pod

import (
	"math"
	"encoding/binary"
	"bytes"
)

//消息
type Message interface {
	//序列化
	Serialize() []byte
	//反序列化
	Deserialize([]byte)
}

type BaseMessage struct {
	MsgType int64
	Length int64
	Data []byte
}

func (sm *BaseMessage) Serialize() []byte {
	b := make([]byte,sm.Length)
	b_buf := bytes.NewBuffer(b[0:8])
	binary.Read(b_buf, binary.BigEndian, sm.MsgType)
	b_buf2 := bytes.NewBuffer(b[8:16])
	binary.Read(b_buf2, binary.BigEndian, sm.Length)
	b[16:sm.Length]=sm.Data
	return b
}

func (sm *BaseMessage) Deserialize(bt []byte) {
	*sm =  BaseMessage{

	}
}


type  HeartbeatMessage struct {

}

func (sm *HeartbeatMessage) Serialize() []byte {
	return []byte(*sm)
}

func (sm *HeartbeatMessage) Deserialize(bt []byte) {
	*sm = StringMsg(bt)
}