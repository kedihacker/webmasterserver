package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	connlist := make([]*websocket.Conn, 0)
	mydialer := websocket.Dialer{}
	for x := 0; x < 1000; x++ {
		conn, _, err := mydialer.Dial("ws://localhost:8080/", nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		connlist = append(connlist, conn)
	}
	log.Println("done")
	// time.Sleep(time.Second * 3)
	for x := 0; x < len(connlist); x++ {
		err := connlist[x].WriteMessage(websocket.TextMessage, []byte("hello"))
		if err != nil {
			log.Fatal(err)
			return
		}
		time.Sleep(time.Millisecond)
		msgtype, msg, err := connlist[x].ReadMessage()
		if err != nil {
			log.Fatal("why", err, x, msgtype, msg)
			return
		}

	}

	log.Println("ended")
	connlist[0].Close()
}
