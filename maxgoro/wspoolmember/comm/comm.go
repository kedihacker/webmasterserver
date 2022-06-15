package comm

import "github.com/gorilla/websocket"

type Commmsg int
type Commstr struct {
	Commmsg Commmsg
	Value   interface{}
}

const (
	Drain    Commmsg = iota
	Newws    Commmsg = iota
	Newframe Commmsg = iota
)

// Contains websocketconn for parsing
type Newwstr struct {
	Wsconn *websocket.Conn
}
