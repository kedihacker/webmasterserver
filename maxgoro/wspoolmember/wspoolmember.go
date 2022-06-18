package wspoolmember

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kedihacker/webmasterserver/maxgoro/wspoolmember/comm"
	"github.com/mailru/easygo/netpoll"
	"github.com/marstr/collection/v2"
)

type Wsmember struct {
	wsconnbuf     *collection.Queue[*timewsconn]
	wsconnbufsize int
	// for reciving ws when its ready no guraantes about being consumed
	rcchan   chan (*websocket.Conn)
	commchan chan (*comm.Commstr)
	//are we accepting new connections?
	Accws        bool
	epool        *netpoll.Poller
	noconcurrent sync.Mutex

	//for general garbage colection
	timey <-chan (time.Time)
}

type timewsconn struct {
	wsconn     *websocket.Conn
	starttime  time.Time
	handledesc *netpoll.Desc
}

func New(size int, rcchan chan (*websocket.Conn)) *Wsmember {
	rb := collection.NewQueue[*timewsconn]()

	tortn := Wsmember{
		wsconnbuf:     rb,
		wsconnbufsize: size,
		rcchan:        rcchan,
		commchan:      make(chan *comm.Commstr, 128),
		Accws:         true,
		noconcurrent:  sync.Mutex{},
	}
	go mainloop(&tortn)
	return &tortn
}

func mainloop(wsm *Wsmember) {
	succes := wsm.noconcurrent.TryLock()
	if !succes {
		log.Println("mainloop: already running")
		return
	}
	tmepepool, err := netpoll.New(&netpoll.Config{})
	wsm.epool = &tmepepool
	if err != nil {
		log.Panic(err)
	}
	wsm.timey = time.Tick(time.Millisecond)
	for {
		select {
		case msg := <-wsm.commchan:
			log.Print("mainloop: got message", msg.Commmsg, " ", len(wsm.rcchan))
			switch msg.Commmsg {
			case comm.Drain:
				wsm.Accws = false
			case comm.Newframe: // ping pong for now
				wsconn := msg.Value.(*timewsconn)
				contenttype, content, err := wsconn.wsconn.ReadMessage()
				if err != nil {
					log.Print("cont read")
					wsconn.wsconn.Close()
					(*(wsm.epool)).Stop(wsconn.handledesc)
				}
				time.Sleep(time.Millisecond)
				err = wsconn.wsconn.WriteMessage(contenttype, content)
				if err != nil {
					log.Println("write error", err)
				}
				err = (*(wsm.epool)).Resume(wsconn.handledesc)

				if err != nil {
					log.Println("resume error", err)
				}
			}
		case <-wsm.timey:
			checkemptyspace(wsm)
			switch {
			case wsm.wsconnbuf.Length() < uint(wsm.wsconnbufsize):
				// fmt.Print("mainloop: accepting new websocket")
				runtime.Gosched()
				select {
				case msg := <-wsm.rcchan:
					newwebssocketfwiend(wsm, msg)
				default:
				}
			default:
			}

			// log.Println("mainloop: default")

		}
	}
}

// gets called when a new websocket is ready to be accepted
// it does labels with current time and add it to epool
func newwebssocketfwiend(wsm *Wsmember, wsconn *websocket.Conn) {
	bufedd := &timewsconn{
		wsconn:     wsconn,
		starttime:  time.Now(),
		handledesc: &netpoll.Desc{},
	}
	wsm.wsconnbuf.Add(bufedd)

	rawdesc, err := netpoll.Handle(bufedd.wsconn.UnderlyingConn(), netpoll.EventRead|netpoll.EventOneShot)

	if err != nil {
		log.Print(err)
		wsconn.Close()
	}
	bufedd.handledesc = rawdesc
	err = (*(wsm.epool)).Start(rawdesc, func(e netpoll.Event) {
		log.Print("new frame")
		wsm.commchan <- &comm.Commstr{
			Commmsg: comm.Newframe,
			Value:   bufedd,
		}
	})
	if err != nil {
		log.Print(err)
		wsconn.Close()
	}

}

func checkemptyspace(wsm *Wsmember) {
	latest, found := wsm.wsconnbuf.Peek()
	if found {
		if time.Since(latest.starttime).Seconds() > 60 {
			log.Println("closing old websocket")

			(*(wsm.epool)).Stop(latest.handledesc)
			wsm.wsconnbuf.Next()
			latest.wsconn.Close()
			checkemptyspace(wsm)
		}

	}
}

func (wsm *Wsmember) Drain(callchan chan (struct{})) {
	wsm.commchan <- &comm.Commstr{
		Commmsg: comm.Drain,
		Value:   nil,
	}

}
