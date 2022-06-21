package main

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"

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

	totalpairs := int64(0)
	wg := sync.WaitGroup{}
	for x := 0; x < 200; x++ {
		wg.Add(1)
		go mamaloop(&wg, &totalpairs)
	}
	wg.Wait()
	log.Println("ended")
	log.Println("total pairs: ", totalpairs)

}

func pairtest(a, b *aaaa, msg Unipubsubmsg) error {
	err := a.Ws.WriteJSON(msg)
	if err != nil {
		return err
	}
	// icmmsg := Unipubsubmsg{}

	_, _, err = b.Ws.ReadMessage()
	if err != nil {
		return err
	}
	return nil
}

func mamaloop(wg *sync.WaitGroup, totalpairs *int64) {
	localparis := 0
	connlist := make([]*aaaa, 0)
	mydialer := websocket.Dialer{}
	for x := 0; x < 100; x++ {
		// log.Print("con num ", x)
		randuuid, _ := uuid.NewRandom()
		conn, _, err := mydialer.Dial("ws://localhost:8080/"+randuuid.String(), nil)
		if err != nil {
			log.Fatal(err)
			return
		}
		connlist = append(connlist, &aaaa{
			Ws: conn,
			Id: randuuid.String(),
		})
	}
	log.Println("done")
	// time.Sleep(time.Second * 3)
	for x := 0; x < 100; x++ {
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
		}
		localparis += 1
	}
	atomic.AddInt64(totalpairs, int64(localparis))
	wg.Done()
}
