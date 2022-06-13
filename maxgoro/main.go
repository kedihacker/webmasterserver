package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mailru/easygo/netpoll"
)

func main() {
	router := mux.NewRouter()
	upy := websocket.Upgrader{}
	epool, err := netpoll.New(&netpoll.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("epool:")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsconn, err := upy.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		rawconn, err := netpoll.Handle(wsconn.UnderlyingConn(), netpoll.EventRead|netpoll.EventOneShot)
		if err != nil {
			log.Println(err)
			return
		}
		epool.Start(rawconn, func(e netpoll.Event) {
			go func() {
				msgtype, msg, err := wsconn.ReadMessage()
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("msgtype:", msgtype, "msg:", string(msg))
				wsconn.WriteMessage(msgtype, msg)
			}()
		})

	})
	//serve touter
	http.ListenAndServe(":8080", router)
}
