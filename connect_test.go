package pod

import (
	"context"
	"fmt"
	"net"
	"testing"
)

func TestNewClientConn(t *testing.T) {
	cc := NewClientConn(0, nil)
	cc.OnActiveTimeout = func(ctx context.Context, conn net.Conn) {
		val := ctx.Value(0)
		fmt.Println(val, "active timeout...")
	}
	fmt.Println(cc)
}
