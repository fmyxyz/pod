package pod

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

//信息类型
type MsgType = byte

func (mt MsgType) String() string {
	switch mt {
	case HEATBEAT:
		return "HEATBEAT"
	case BINARY:
		return "BINARY"
	case BIGBINARY:
		return "BIGBINARY"
	case TEXT:
		return "TEXT"
	case BIGTEXT:
		return "BIGTEXT"
	case MORETEXT:
		return "MORETEXT"
	case MOREBINARY:
		return "MOREBINARY"
	default:
		return "Unknow"
	}
}

const (
	_ MsgType = iota
	//心跳
	HEATBEAT
	//二进制数据
	BINARY
	//大二进制数据
	//先传输类型大小等mate信息
	BIGBINARY
	//文本数据
	TEXT
	//大文本数据
	//先传输类型大小等mate信息
	BIGTEXT
	//后续还有文本数据
	MORETEXT
	//后续还有二进制数据
	MOREBINARY
	//消息序列sequence
	SEQUENCE
	//其他
)

//元数据字节数
const METADATALENGTH int = 9

//可读写消息
//读写顺序为 消息类型,字节长度,数据
type Message interface {
	//序列化 包括 元数据
	Serialize(writer io.Writer) error
	//反序列化 具体不包括 元数据
	Deserialize(reader io.Reader) error
	//消息类型
	MessageType() MsgType
	//字节长度
	Size() int64
}

//消息元数据
type Metadata struct {
	MsgType MsgType
	Length  int64
}

func (h Metadata) Size() int64 {
	return h.Length
}

func (h Metadata) MessageType() MsgType {
	return h.MsgType
}

func (h *Metadata) Serialize(writer io.Writer) error {
	var err error
	bs := make([]byte, METADATALENGTH)
	bs[0] = h.MsgType
	b_buf := bytes.NewBuffer(bs[1:METADATALENGTH])
	binary.Write(b_buf, binary.BigEndian, h.Length)
	_, err = writer.Write(bs)
	if err != nil {
		return err
	}
	return err
}

func (h *Metadata) Deserialize(reader io.Reader) error {
	var err error
	bt := make([]byte, METADATALENGTH)
	_, err = reader.Read(bt)
	if err != nil {
		return err
	}
	h.MsgType = MsgType(bt[0])
	b_buf := bytes.NewBuffer(bt[1:METADATALENGTH])
	err = binary.Read(b_buf, binary.BigEndian, &h.Length)
	return err
}


// 注册消息
type RegistryedMessage struct {
	registryedMessage map [MsgType] Message
	metadata Metadata
}

func (h RegistryedMessage) Size() int64 {
	return h.registryedMessage[h.metadata.MsgType].Size()
}

func (h RegistryedMessage) MessageType() MsgType {
	return  h.metadata.MsgType
}

func (h *RegistryedMessage) Serialize(writer io.Writer) error {
	var err error
	msg:= h.registryedMessage[h.metadata.MsgType]
	err=msg.Serialize(writer)
	return err
}

func (h *RegistryedMessage) Deserialize(reader io.Reader) error {
	var err error
	msg:= h.registryedMessage[h.metadata.MsgType]
	err=msg.Deserialize(reader)
	return err
}



// 一系列消息
type MessageQueue []Message

func (h MessageQueue) Size() int64 {
	return -1
}

func (h MessageQueue) MessageType() MsgType {
	return SEQUENCE
}

func (h *MessageQueue) Serialize(writer io.Writer) error {
	var err error
	for _, v := range *h {
		err = v.Serialize(writer);
		if err != nil {
			return err
		}
	}
	return err
}

func (h *MessageQueue) Deserialize(reader io.Reader) error {
	var err error
	for _, v := range *h {
		err = v.Deserialize(reader);
		if err != nil {
			return err
		}
	}
	return err
}

//心跳数据
type HeatbeatMsg struct {
	Metadata
	Duration time.Duration
	//data io.ReadWriter
}

//默认30s
func NewHeatbeatMsg() HeatbeatMsg {
	return HeatbeatMsg{Duration: 30 * time.Second}
}

func (h HeatbeatMsg) Size() int64 {
	return h.Length
}

func (h HeatbeatMsg) MessageType() MsgType {
	return h.MsgType
}

func (h *HeatbeatMsg) Serialize(writer io.Writer) error {
	var err error
	//消息长度为8字节
	h.Length=8
	h.MsgType=HEATBEAT
	err = h.Metadata.Serialize(writer)
	if err != nil {
		return err
	}
	bt := make([]byte, 8)
	b_buf := bytes.NewBuffer(bt)
	binary.Write(b_buf, binary.BigEndian, h.Duration)
	_, err = writer.Write(b_buf.Bytes())
	return err
}

func (h *HeatbeatMsg) Deserialize(reader io.Reader) error {
	var err error
	bt := make([]byte, 8)
	b_buf := bytes.NewBuffer(bt)
	err = binary.Read(b_buf, binary.BigEndian, &h.Duration)
	_, err = reader.Read(bt)
	return err
}
