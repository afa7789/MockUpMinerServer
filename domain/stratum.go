package domain

import "net"

// Stratum Entry which we receive
type StratumEntryRPC struct {
	ID     int64         `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type StratumManager interface {
	HandleConn(c net.Conn)
}

// type StratumSubscribeReturn struct {
// 	ID    int          `json:"id"`
// 	Error StratumError `json:"error"`
// }

type StratumError struct {
	ID           int         `json:"id"`
	ErrorMessage string      `json:"error"`
	Traceback    interface{} `json:"traceback"` // didn't understood what a traceback is.
}
