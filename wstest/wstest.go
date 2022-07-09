package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type aaaa struct {
	Ws *websocket.Conn
	Id string
}

type Unipubsubmsg struct {
	Keyy  string `json:"key"`
	Value string `json:"value"`
}

func main() {
	log.Println("start")

	starttime := time.Now()

	totalpairs := int64(0)
	wg := sync.WaitGroup{}
	for x := 0; x < 1; x++ {
		wg.Add(1)
		go mamaloop(&wg, &totalpairs)
	}
	wg.Wait()
	log.Println("ended")
	log.Println("total pairs: ", totalpairs)
	log.Println(time.Since(starttime))

}

func pairtest(a, b *aaaa, msg Unipubsubmsg) error {
	log.Println("pairtest")
	err := a.Ws.WriteJSON(msg)
	if err != nil {
		return err
	}
	// icmmsg := Unipubsubmsg{}
	b.Ws.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
	_, _, err = b.Ws.ReadMessage()
	if err != nil {
		return errors.New(err.Error() + " " + b.Id + " " + a.Id)
	}

	return nil
}

func mamaloop(wg *sync.WaitGroup, totalpairs *int64) {
	localparis := 0
	connlist := make([]*aaaa, 0)

	for x := 0; x < 16; x++ {
		log.Print("con num ", x)
		randuuid, _ := uuid.NewRandom()
		conn, httpres, err := websocket.DefaultDialer.Dial("wss://yenicericopybackend.herokuapp.com/ "+randuuid.String(), nil)

		if err != nil && conn != nil && httpres.StatusCode != 101 {
			log.Fatal(err)
			return
		}
		connlist = append(connlist, &aaaa{
			Ws: conn,
			Id: randuuid.String(),
		})
	}
	time.Sleep(time.Second * 10)
	log.Println("done")
	// time.Sleep(time.Second * 3)
	for x := 0; x < 64; x++ {
		firstcandidateindex := rand.Intn(len(connlist))
		first := connlist[firstcandidateindex]
		secondcandidate := rand.Intn(len(connlist) - 1)
		if secondcandidate >= firstcandidateindex {
			secondcandidate++
		}

		second := connlist[secondcandidate]
		if first.Id == second.Id {
			log.Println("same")
			continue
		}
		// log.Print("pairing ", x)
		err := pairtest(first, second, Unipubsubmsg{
			Keyy:  second.Id,
			Value: uuid.New().String(),
		})
		if err != nil {
			log.Println(err)
			localparis--
		}
		localparis += 1
	}
	atomic.AddInt64(totalpairs, int64(localparis))
	wg.Done()
}
