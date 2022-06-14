package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	connlist := make([]*websocket.Conn, 0)
	mydialer := websocket.Dialer{}
	for x := 0; x < 50000; x++ {
		log.Print("con num ", x)
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
		msgtype, msg, err := connlist[x].ReadMessage()
		log.Println(x)
		if err != nil {
			log.Fatal("why", err, x, msgtype, msg)
			return
		}

	}
	for x := 0; x < len(connlist); x++ {
		err := connlist[x].WriteMessage(websocket.TextMessage, []byte("hello"))
		if err != nil {
			log.Fatal(err)
			return
		}
		msgtype, msg, err := connlist[x].ReadMessage()
		log.Println(x)
		if err != nil {
			log.Fatal("why", err, x, msgtype, msg)
			return
		}
		// connlist[x].Close()

	}

	log.Println("ended")
	connlist[0].Close()
}
