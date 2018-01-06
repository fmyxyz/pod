package utils

import "sync"

var id int64 = 0

var mux sync.Mutex

func init() {
	mux = sync.Mutex{}
}

func GenerationId() int64 {
	mux.Lock()
	defer mux.Unlock()
	id += 1
	return id
}
