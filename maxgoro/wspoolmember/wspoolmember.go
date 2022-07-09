package wspoolmember

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kedihacker/webmasterserver/maxgoro/msgrouter"
	"github.com/kedihacker/webmasterserver/maxgoro/wspoolmember/comm"
	"github.com/mailru/easygo/netpoll"
	"github.com/marstr/collection/v2"
)

type Wsmember struct {
	wsconnbuf      *collection.Queue[*Timewsconn]
	wsconnbufsize  int
	wsconnlifetime float64
	// for reciving ws when its ready no guraantes about being consumed
	rcchan   chan (comm.Newfwiendmetadata)
	commchan chan (*comm.Commstr)
	//are we accepting new connections?
	Accws bool
	epool *netpoll.Poller

	// for preventng accidental launcihng of multiple goroutines
	noconcurrent sync.Mutex
	pubsunpool   *msgrouter.Msgrouter[comm.Unipubsubmsg]
	// for getting mesages from pubsub
	punsubinc chan (comm.Unipubsubmsg)
	//for general garbage colection
	timey         <-chan (time.Time)
	wsmap         map[string]*Timewsconn
	casulydeleted map[string]struct{}
}

type Timewsconn struct {
	id         string
	wsconn     *websocket.Conn
	starttime  time.Time
	handledesc *netpoll.Desc
	e          netpoll.Event
}

func New(size int, rcchan chan (comm.Newfwiendmetadata), pubsubpool *msgrouter.Msgrouter[comm.Unipubsubmsg], lifetime float64) *Wsmember {
	rb := collection.NewQueue[*Timewsconn]()

	tortn := Wsmember{
		wsconnlifetime: lifetime,
		wsconnbuf:      rb,
		wsconnbufsize:  size,
		rcchan:         rcchan,
		commchan:       make(chan *comm.Commstr, 128),
		Accws:          true,
		noconcurrent:   sync.Mutex{},
		punsubinc:      make(chan (comm.Unipubsubmsg), 128),
		pubsunpool:     pubsubpool,
		wsmap:          make(map[string]*Timewsconn),
		casulydeleted:  make(map[string]struct{}),
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
	brandnewticker := time.NewTicker(time.Millisecond)
	wsm.timey = brandnewticker.C
	defer brandnewticker.Stop()
	for {
		select {
		case msg := <-wsm.commchan:
			// log.Println(len(wsm.commchan))
			// log.Print("mainloop: got message", msg.Commmsg, " ", len(wsm.rcchan))
			switch msg.Commmsg {
			case comm.Drain:
				wsm.Accws = false
			case comm.Newframe:
				wsconn := msg.Value.(*Timewsconn)
				// log.Println("event type :", ((wsconn.e & (netpoll.EventHup | netpoll.EventReadHup | netpoll.EventWriteHup)) != 0))
				// log.Println("event type :", wsconn.e)
				if (wsconn.e & (netpoll.EventHup | netpoll.EventReadHup | netpoll.EventWriteHup)) != 0 {
					log.Println("mainloop: websocket closed")
					delete(wsm.wsmap, wsconn.id)
					wsm.pubsunpool.Unsub(wsconn.id, wsm.punsubinc)
					(*(wsm.epool)).Stop(wsconn.handledesc)
					wsconn.wsconn.Close()
					continue
				}

				rcvmsg := comm.Unipubsubmsg{}
				err := wsconn.wsconn.ReadJSON(&rcvmsg)
				if err != nil {
					log.Print("can't read", err)
					wsconn.wsconn.Close()
					delete(wsm.wsmap, wsconn.id)
					wsm.pubsunpool.Unsub(wsconn.id, wsm.punsubinc)
					(*(wsm.epool)).Stop(wsconn.handledesc)
					continue
				}
				// log.Print("mainloop: sent message", rcvmsg)
				wsm.pubsunpool.Pub(rcvmsg.Keyy, rcvmsg)
				err = (*(wsm.epool)).Resume(wsconn.handledesc)
				if err != nil {
					wsconn.wsconn.Close()
					delete(wsm.wsmap, wsconn.id)
					wsm.pubsunpool.Unsub(wsconn.id, wsm.punsubinc)
					log.Println("resume error", err)
					(*(wsm.epool)).Stop(wsconn.handledesc)
					continue
				}
			}
		case msg := <-wsm.punsubinc:
			// log.Print("mainloop: got message ", msg.Keyy, " ", len(wsm.rcchan))
			if _, found := wsm.casulydeleted[msg.Keyy]; wsm.wsmap[msg.Keyy] == nil && found {
				continue
			}
			if wsm.wsmap[msg.Keyy].wsconn == nil {
				log.Panic("nowsconn")
			}
			err := wsm.wsmap[msg.Keyy].wsconn.WriteJSON(msg)
			if err != nil {
				log.Println("error writing to websocket", err)
				(*(wsm.epool)).Stop(wsm.wsmap[msg.Keyy].handledesc)
				wsm.wsmap[msg.Keyy].wsconn.Close()
				wsm.pubsunpool.Unsub(msg.Keyy, wsm.punsubinc)
				delete(wsm.wsmap, msg.Keyy)
			}
		case <-wsm.timey:
			checkemptyspace(wsm)

			if wsm.Accws {
				switch {
				case wsm.wsconnbuf.Length() < uint(wsm.wsconnbufsize):
					// fmt.Print("mainloop: accepting new websocket")
					select {
					case msg := <-wsm.rcchan:
						newwebssocketfwiend(wsm, msg)
					default:
					}
				default:
				} // log.Println("mainloop: default")
			} else {
				if wsm.wsconnbuf.IsEmpty() {
					return
				}
			}
		}
	}
}

// gets called when a new websocket is ready to be accepted
// it does labels with current time and add it to epool
func newwebssocketfwiend(wsm *Wsmember, wsconn comm.Newfwiendmetadata) {

	bufedd := &Timewsconn{
		id:         wsconn.Kv["listenker"],
		wsconn:     wsconn.Conn,
		starttime:  time.Now(),
		handledesc: &netpoll.Desc{},
	}
	wsm.wsconnbuf.Add(bufedd)

	rawdesc, err := netpoll.Handle(bufedd.wsconn.UnderlyingConn(), netpoll.EventRead|netpoll.EventOneShot)

	if err != nil {
		log.Print(err)
		wsconn.Conn.Close()
	}
	bufedd.handledesc = rawdesc
	err = (*(wsm.epool)).Start(rawdesc, func(e netpoll.Event) {
		// log.Print("new frame")
		bufedd.e = e
		wsm.commchan <- &comm.Commstr{
			Commmsg: comm.Newframe,
			Value:   bufedd,
		}
	})
	if err != nil {
		log.Print("cant register", err)
		wsconn.Conn.Close()
	}
	wsm.wsmap[wsconn.Kv["listenker"]] = bufedd
	wsm.pubsunpool.Sub(wsconn.Kv["listenker"], wsm.punsubinc)

}

func checkemptyspace(wsm *Wsmember) {
	latest, found := wsm.wsconnbuf.Peek()
	if found {
		if time.Since(latest.starttime).Seconds() > wsm.wsconnlifetime {
			log.Println("closing old websocket")

			(*(wsm.epool)).Stop(latest.handledesc)
			wsm.wsconnbuf.Next()
			latest.wsconn.Close()
			wsm.pubsunpool.Unsub(latest.id, wsm.punsubinc)
			delete(wsm.wsmap, latest.id)
			wsm.casulydeleted[latest.id] = struct{}{}
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
