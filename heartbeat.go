package pod

import (
	"log"

	"os"
)

var conns map[int64]*ConnAction
var logger *log.Logger

func AddConn(nid int64, ca *ConnAction) {
	conns[nid] = ca
}

func init() {
	conns = make(map[int64]*ConnAction)
	logger=log.New(os.Stdout,"heartbeat",log.LstdFlags|log.Llongfile)
	go heartbeat(conns)
}

func heartbeat(conns map[int64]*ConnAction)  {
	for{
		for k,v:= range conns  {
			logger.Println("conn id :",k)
			if (*v).GetTimeout()==0{
				(*v).PullMessage(func() int64 {

				},)
			}
		}
	}
}