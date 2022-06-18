package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kedihacker/webmasterserver/maxgoro/wspoolmember"
)

type Notifier[T any] interface {
	GetMessagechan(string) chan (T)
	AddMessage(string, T) error

	// This function is concurrently and blokingly called
	AddNotifierfunc(func(T)) error
	DeleteNotifier()
}

var incws = make(chan (*websocket.Conn), 512)

func main() {
	router := mux.NewRouter()
	upy := websocket.Upgrader{
		WriteBufferPool: &sync.Pool{},
	}
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "127.18.0.1:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// epool, err := netpoll.New(&netpoll.Config{})
	// if err != nil {
	// 	panic(err)
	// }

	for x := 0; x < 256; x++ {
		wspoolmember.New(64, incws)
	}
	log.Println("epool:")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsconn, err := upy.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			wsconn.Close()
			return
		}
		incws <- wsconn
		// rawconn, err := netpoll.Handle(wsconn.UnderlyingConn(), netpoll.EventRead|netpoll.EventOneShot)
		// if err != nil {
		// 	rawconn.Close()
		// 	wsconn.Close()
		// 	log.Println(err)
		// 	return
		// }
		// epool.Start(rawconn, func(e netpoll.Event) {

		// 	go handleit(wsconn, epool, rawconn)
		// })

	})

	//serve touter
	http.ListenAndServe(":8080", router)
}

// func handleit(wsconn *websocket.Conn, epool netpoll.Poller, rawconn *netpoll.Desc) {

// 	msgtype, msg, err := wsconn.ReadMessage()
// 	if err != nil {
// 		log.Println(err)

// 		epool.Stop(rawconn)
// 		rawconn.Close()
// 		wsconn.Close()
// 		return
// 	}
// 	time.Sleep(time.Millisecond)
// 	log.Println("msgtype:", msgtype, "msg:", string(msg))
// 	wsconn.WriteMessage(msgtype, msg)
// 	rawconn, err = netpoll.Handle(wsconn.UnderlyingConn(), netpoll.EventRead|netpoll.EventOneShot)
// 	if err != nil {
// 		log.Println(err)
// 		epool.Stop(rawconn)
// 		rawconn.Close()
// 		wsconn.Close()
// 		return
// 	}
// 	epool.Start(rawconn, func(e netpoll.Event) {
// 		go handleit(wsconn, epool, rawconn)
// 	})
// }
