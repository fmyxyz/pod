package pod

import (
//"context"
)

type ClientConn struct {
	BaseConn
}

var hbm HeatbeatMsg = NewHeatbeatMsg()

func NewClientConn(netid int64) *ClientConn {
	cc := &ClientConn{}
	//cc.OnActiveTimeout = func(ctx context.Context, conn ConnAction) {
	//	select {
	//	case <-ctx.Done():
	//		var message Message = &hbm
	//		conn.PullMessage(message)
	//	}
	//}
	return cc
}
