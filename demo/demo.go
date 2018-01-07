package main

import (
	"encoding/binary"
	"bytes"
	"fmt"
)

func main() {
	b_buf := bytes.NewBuffer([]byte{})
	var i int64 = 0xffffaa88
	binary.Write(b_buf, binary.BigEndian, i)
	fmt.Println(b_buf.Bytes())
	var b=b_buf.Bytes()
	b2 := []byte{8: 0}
	a :=copy(b2[1:9], b)
	fmt.Println(b2,a)

}
