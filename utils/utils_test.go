package utils

import "testing"

func TestGenerationId(t *testing.T) {
	ids := make(chan int64, 1048576)
	for i := 0; i < 1048576; i++ {
		go func() {
			id := GenerationId()
			ids <- id
		}()
	}
	var id int64
	for i := 0; i < 1048576; i++ {
		select {
		case id = <-ids:
			if i != int(id) {
				println(i, id)
			}
		}
	}
}
