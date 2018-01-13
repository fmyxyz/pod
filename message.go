package pod

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
	"errors"
	"fmt"
	"log"
)

//信息类型
type MsgType = byte

func MsgTypeString(mt MsgType) string {
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
	//获取消息类型
	GetMetadata() Metadata
	//设置消息类型
	SetMetadata(md Metadata)
}

//消息元数据
type Metadata struct {
	MsgType MsgType
	Length  int64
}

func (h Metadata) GetMetadata() Metadata {
	return h
}

func (h *Metadata) SetMetadata(md Metadata) {
	h.Length = md.Length
	h.MsgType = md.MsgType
}

func (h *Metadata) Serialize(writer io.Writer) error {
	var err error
	bs := make([]byte, METADATALENGTH)
	bs[0] = h.MsgType
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, h.Length)
	lbs := b_buf.Bytes(); //h.Length 二进制数据
	copy(bs[1:METADATALENGTH], lbs)
	_, err = writer.Write(bs)
	if err != nil {
		return err
	}
	return err
}

func (h Metadata) String() string {
	str := fmt.Sprintln("message type:", MsgTypeString(h.MsgType), ";", "message length:", h.Length)
	return str
}

func (h *Metadata) Deserialize(reader io.Reader) error {
	var err error
	var l int
	bt := make([]byte, METADATALENGTH)
	l, err = reader.Read(bt)
	if err != nil || l != METADATALENGTH {
		return err
	}
	h.MsgType = MsgType(bt[0])
	b_buf := bytes.NewBuffer(bt[1:METADATALENGTH])
	err = binary.Read(b_buf, binary.BigEndian, &h.Length)
	return err
}

// 注册消息
type RegistryedMessage struct {
	registryedMessage map[MsgType]Message
	metadata          Metadata
}

func RegistryMessage(mt MsgType, msg Message) *RegistryedMessage {
	rm := new(RegistryedMessage)
	rm.registryedMessage = make(map[MsgType]Message)
	rm.registryedMessage[mt] = msg
	return rm
}

func (rm *RegistryedMessage) RegistryMessage(mt MsgType, msg Message) *RegistryedMessage {
	rm.registryedMessage[mt] = msg
	return rm
}

func (h *RegistryedMessage) SetMetadata(md Metadata) {
	h.metadata = md
}

func (h RegistryedMessage) GetMetadata() Metadata {
	return h.metadata
}

func (h *RegistryedMessage) Serialize(writer io.Writer) error {
	var err error
	msg := h.registryedMessage[h.metadata.MsgType]
	//设置元数据
	msg.SetMetadata(h.metadata)
	err = msg.Serialize(writer)
	return err
}

func (h *RegistryedMessage) Deserialize(reader io.Reader) error {
	var err error
	msg := h.registryedMessage[h.metadata.MsgType]
	if msg == nil {
		return errors.New(fmt.Sprint("未注册：", h.metadata.MsgType))
	}
	//设置元数据
	msg.SetMetadata(h.metadata)
	err = msg.Deserialize(reader)
	return err
}

// 一系列消息
type MessageQueue []Message

func (h MessageQueue) GetMetadata() Metadata {
	for _, v := range h {
		return v.GetMetadata()
	}
	var md Metadata
	return md
}

func (h *MessageQueue) SetMetadata(md Metadata) {
	for _, v := range *h {
		v.SetMetadata(md)
	}
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
}

//默认30s
func NewHeatbeatMsg() HeatbeatMsg {
	hm := HeatbeatMsg{Duration: 30 * time.Second}
	hm.MsgType = HEATBEAT
	return hm
}

func (h *HeatbeatMsg) SetMetadata(md Metadata) {
	h.Metadata = md
}

func (h HeatbeatMsg) GetMetadata() Metadata {
	return h.Metadata
}

func (h *HeatbeatMsg) Serialize(writer io.Writer) error {
	var err error
	//消息长度为8字节
	h.Length = 8
	h.MsgType = HEATBEAT
	err = h.Metadata.Serialize(writer)
	if err != nil {
		return err
	}
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, h.Duration)
	bt := b_buf.Bytes()
	_, err = writer.Write(bt)
	return err
}

func (h *HeatbeatMsg) Deserialize(reader io.Reader) error {
	var err error
	bt := make([]byte, 8)
	b_buf := bytes.NewBuffer(bt)
	_, err = reader.Read(bt)
	err = binary.Read(b_buf, binary.BigEndian, &h.Duration)
	return err
}

type BinaryMsg struct {
	Metadata
	Data []byte
}

func NewBinaryMsg() BinaryMsg {
	bm := BinaryMsg{}
	bm.MsgType = BINARY
	return bm
}

func (h *BinaryMsg) SetMetadata(md Metadata) {
	h.Metadata = md
}

func (h BinaryMsg) GetMetadata() Metadata {
	return h.Metadata
}

func (h *BinaryMsg) Serialize(writer io.Writer) error {
	var err error
	//消息长度为8字节
	if h.Data != nil {
		h.Length = int64(len(h.Data))
	}

	err = h.Metadata.Serialize(writer)
	if err != nil {
		return err
	}
	if h.Data != nil {
		_, err = writer.Write(h.Data)
	}
	return err
}

func (h *BinaryMsg) Deserialize(reader io.Reader) error {
	var err error
	h.Data = make([]byte, h.Length)
	_, err = reader.Read(h.Data)
	log.Println(string(h.Data))
	return err
}

type TextMsg struct {
	Metadata
	Data string
}

func NewTextMsg() TextMsg {
	bm := TextMsg{}
	bm.MsgType = BINARY
	return bm
}

func (h *TextMsg) SetMetadata(md Metadata) {
	h.Metadata = md
}

func (h TextMsg) GetMetadata() Metadata {
	return h.Metadata
}

func (h *TextMsg) Serialize(writer io.Writer) error {
	var err error
	h.Length = int64(len(h.Data))
	err = h.Metadata.Serialize(writer)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(h.Data))
	return err
}

func (h *TextMsg) Deserialize(reader io.Reader) error {
	var err error
	var data = make([]byte, h.Length)
	_, err = reader.Read(data)
	h.Data = string(data)
	log.Println(string(h.Data))
	return err
}
