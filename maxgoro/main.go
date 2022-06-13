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

	if err != nil {
		panic(err)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsconn, err := upy.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

	})
}
