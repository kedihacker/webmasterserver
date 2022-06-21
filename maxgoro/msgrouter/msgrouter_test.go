package msgrouter_test

import (
	"testing"

	"github.com/kedihacker/webmasterserver/maxgoro/msgrouter"
)

func TestMsgrouter(t *testing.T) {
	// TODO
	testinst := msgrouter.New[string]()
	firstchan := make(chan string, 2)
	secondchan := make(chan string, 2)
	testinst.Sub("sendfirst", firstchan)
	testinst.Sub("sendsecond", secondchan)
	testinst.Pub("sendfirst", "hifirst")
	testinst.Pub("sendsecond", "hisecond")
	if <-firstchan != "hifirst" {
		t.Error("first chan not working")
	}
	if <-secondchan != "hisecond" {
		t.Error("second chan not working")
	}
}
