package pod

import (
	"encoding/binary"
	"bytes"
	"io"
	"time"
)

//可读写消息
type Message interface {
	//序列化
	Serialize(writer io.Writer)
	//反序列化
	Deserialize(reader io.Reader)
}

type heatbeatMsg struct {
	MsgType int64
	Length int64
	Duration time.Duration
	//data io.ReadWriter
}

func  NewHeatbeatMsg() heatbeatMsg  {
	return heatbeatMsg{MsgType:1,Length:24,Duration:30*time.Second}
}

func (h *heatbeatMsg)Serialize(writer io.Writer) {
	bt:=make([]byte,8)

	b_buf:=bytes.NewBuffer(bt)
	binary.Write(b_buf,binary.BigEndian,h.MsgType)
	writer.Write(b_buf.Bytes())

	b_buf=bytes.NewBuffer(bt)
	binary.Write(b_buf,binary.BigEndian,h.Length)
	writer.Write(b_buf.Bytes())

	b_buf=bytes.NewBuffer(bt)
	binary.Write(b_buf,binary.BigEndian,h.Duration)
	writer.Write(b_buf.Bytes())
}

func (h *heatbeatMsg)Deserialize(reader io.Reader){
	bt:=make([]byte,8)
	reader.Read(bt)
	b_buf:=bytes.NewBuffer(bt)
	binary.Read(b_buf,binary.BigEndian,h.MsgType)

	reader.Read(bt)
	b_buf=bytes.NewBuffer(bt)
	binary.Read(b_buf,binary.BigEndian,h.Length)

	reader.Read(bt)
	b_buf=bytes.NewBuffer(bt)
	binary.Read(b_buf,binary.BigEndian,h.Duration)
}
