package pod

type ServerConn struct {
	BaseConn
}

func NewServerConn(netid int64) *ServerConn {
	cc := &ServerConn{}
	return cc
}
