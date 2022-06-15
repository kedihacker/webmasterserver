package wspoolmember

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kedihacker/webmasterserver/maxgoro/wspoolmember/comm"
	"github.com/mailru/easygo/netpoll"
	"github.com/marstr/collection/v2"
)

type Wsmember struct {
	wsconnbuf *collection.Queue[*timewsconn]
	rcchan    chan (websocket.Conn)
	commchan  chan (comm.Commstr)
	Accws     bool
	epool     netpoll.Poller
}

type timewsconn struct {
	wsconn     *websocket.Conn
	starttime  time.Time
	handledesc []*netpoll.Desc
}

func New(size int, rcchan chan (websocket.Conn)) *Wsmember {
	rb := collection.NewQueue[*timewsconn]()

	tortn := Wsmember{
		wsconnbuf: rb,
		rcchan:    rcchan,
		commchan:  make(chan comm.Commstr, 128),
		Accws:     false,
	}
	go mainloop(&tortn)
	return &tortn
}

func mainloop(wsm *Wsmember) {
	tmepepool, err := netpoll.New(&netpoll.Config{})
	wsm.epool = tmepepool
	if err != nil {
		log.Panic(err)
	}
	for {
		msg, ok := <-wsm.commchan
		if ok {
			switch msg.Commmsg {
			case comm.Drain:
				wsm.Accws = false
			case comm.Newws:
				wsstr := msg.Value.(comm.Newwstr).Wsconn
				if wsm.Accws {
					if wsm.wsconnbuf.IsEmpty() {
						bufedd := &timewsconn{
							wsconn:     wsstr,
							starttime:  time.Now(),
							handledesc: []*netpoll.Desc{},
						}
						wsm.wsconnbuf.Add(bufedd)

						rawdesc, err := netpoll.Handle(wsstr.UnderlyingConn(), netpoll.EventRead)

						if err != nil {
							log.Print(err)
							wsstr.Close()
						}
						bufedd.handledesc = append(bufedd.handledesc, rawdesc)
						wsm.epool.Start(rawdesc, func(e netpoll.Event) {
							wsm.commchan <- comm.Commstr{
								Commmsg: comm.Newframe,
								Value:   bufedd,
							}
						})
					} else {
						checkemptyspace(wsm)
						bufedd := &timewsconn{
							wsconn:     wsstr,
							starttime:  time.Now(),
							handledesc: []*netpoll.Desc{},
						}
						wsm.wsconnbuf.Add(bufedd)

						rawdesc, err := netpoll.Handle(wsstr.UnderlyingConn(), netpoll.EventRead|netpoll.EventOneShot)

						if err != nil {
							log.Print(err)
							wsstr.Close()
						}
						bufedd.handledesc = append(bufedd.handledesc, rawdesc)
						wsm.epool.Start(rawdesc, func(e netpoll.Event) {
							wsm.epool.Resume(rawdesc)
							wsm.commchan <- comm.Commstr{
								Commmsg: comm.Newframe,
								Value:   bufedd,
							}
						})
					}
				} else {
					log.Print("sent a wsconn to draining wspoolmember")
					msg.Value.(comm.Newwstr).Wsconn.Close()
				}
			case comm.Newframe:
				wsconn := msg.Value.(*timewsconn)
				contenttype, content, err := wsconn.wsconn.ReadMessage()
				if err != nil {
					wsconn.wsconn.Close()
					for _, x := range wsconn.handledesc {
						wsm.epool.Stop(x)
					}

				}
				wsconn.wsconn.WriteMessage(contenttype, content)
			}
		} else {
			checkemptyspace(wsm)
		}

	}
}

func checkemptyspace(wsm *Wsmember) {
	latest, found := wsm.wsconnbuf.Peek()
	if found {
		if time.Since(latest.starttime).Seconds() > 60 {
			for _, x := range latest.handledesc {
				wsm.epool.Stop(x)
			}
			latest.wsconn.Close()
			checkemptyspace(wsm)
		}

	}
}

func (wsm *Wsmember) Drain(callchan chan (struct{})) {
	wsm.commchan <- comm.Commstr{
		Commmsg: comm.Drain,
		Value:   nil,
	}

}
