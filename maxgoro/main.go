package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kedihacker/webmasterserver/maxgoro/msgrouter"
	"github.com/kedihacker/webmasterserver/maxgoro/wspoolmember"
	"github.com/kedihacker/webmasterserver/maxgoro/wspoolmember/comm"
)

type Notifier[T any] interface {
	GetMessagechan(string) chan (T)
	AddMessage(string, T) error

	// This function is concurrently and blokingly called
	AddNotifierfunc(func(T)) error
	DeleteNotifier()
}

var incws = make(chan (comm.Newfwiendmetadata), 512)

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
	pubsubpool := msgrouter.New[comm.Unipubsubmsg]()
	for x := 0; x < 512; x++ {
		wspoolmember.New(512, incws, pubsubpool, 60)
	}
	log.Println("epool:")
	router.HandleFunc("/{listenker:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$}", func(w http.ResponseWriter, r *http.Request) {
		kv := mux.Vars(r)
		wsconn, err := upy.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			wsconn.Close()
			return
		}
		incws <- comm.Newfwiendmetadata{
			Conn: wsconn,
			Kv:   kv,
		}
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

	//serve router
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Panic(err)
	}
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
