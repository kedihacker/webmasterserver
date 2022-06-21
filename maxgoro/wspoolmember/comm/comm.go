package comm

import "github.com/gorilla/websocket"

type Commmsg int
type Commstr struct {
	Commmsg Commmsg
	Value   interface{}
}

const (
	Drain Commmsg = iota
	Newws Commmsg = iota
	// there is new data on the websocket
	Newframe Commmsg = iota
)

type Newfwiendmetadata struct {
	Conn *websocket.Conn
	Kv   map[string]string
}

// Contains websocketconn for parsing
type Newwstr struct {
	Wsconn     *websocket.Conn
	Returnchan chan (*websocket.Conn)
}

type Unipubsubmsg struct {
	Keyy  string `json:"key"`
	Value string `json:"value"`
}

func (u *Unipubsubmsg) GettKeyy() string {
	return u.Keyy
}
