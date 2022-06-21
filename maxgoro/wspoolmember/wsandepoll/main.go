package wsandepoll

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/kedihacker/webmasterserver/maxgoro/wspoolmember"
	"github.com/mailru/easygo/netpoll"
)

type Pollandws struct {
	Poller *netpoll.Poller
	Ws     *websocket.Conn
	raw    *netpoll.Desc
	Expiry time.Time
}

func New(poll *netpoll.Poller, ws *websocket.Conn, srtfunc func(e netpoll.Event), flags netpoll.Event, t time.Time) (*Pollandws, error) {
	raw, err := netpoll.Handle(ws.UnderlyingConn(), flags)
	if err != nil {
		return nil, err
	}
	err = (*(poll)).Start(raw, srtfunc)
	if err != nil {
		return nil, err
	}
	return &Pollandws{
		Poller: poll,
		Ws:     ws,
		Expiry: t,
		raw:    raw,
	}, nil
}

func (p *Pollandws) Resume() error {
	err := (*(p.Poller)).Resume(p.raw)
	if err != nil {
		return err
	}
	return nil
}

func (p *Pollandws) Close() error {
	err := (*(p.Poller)).Stop(p.raw)
	if err != nil {
		return err
	}
	err = p.Ws.Close()
	if err != nil {
		return err
	}
	return nil
}

type Mapcordinator struct {
	Map   map[string]*wspoolmember.Timewsconn
	epoll *netpoll.Poller
}
