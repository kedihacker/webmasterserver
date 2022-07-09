package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {
	randuuid, _ := uuid.NewRandom()
	conn, _, err := websocket.DefaultDialer.Dial("https://yenicericopybackend.herokuapp.com/"+randuuid.String(), nil)
	if err != nil {
		panic(err)
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		panic(err)
	}
}
