package remote_poller

import (
	"testing"
)

type testReceiver struct {
	notify      chan<- bool
	delReceived bool
	modReceived bool
}

func (r testReceiver) OnFileAdded(e Event) {
	r.notify <- true
}

func (r testReceiver) OnFileDeleted(e Event) {
	r.notify <- true
}
func (r testReceiver) OnFileModified(e Event) {
	r.notify <- true
}

func TestEventTrigger_OnEvents(t *testing.T) {

	listeners := make([]Receiver, 0)

	addChan, modChan, delChan := make(chan Event), make(chan Event), make(chan Event)
	notifyChan := make(chan bool)

	r := testReceiver{notify: notifyChan}

	listeners = append(listeners, r, r)

	em := eventTriggerManager{receivers: listeners}

	testElement := &testElement{name: "testElement"}

	go em.OnFileAdded(addChan)
	go em.OnFileModified(modChan)
	go em.OnFileDeleted(delChan)

	addChan <- &triggeredEvent{e: testElement}
	modChan <- &triggeredEvent{e: testElement}
	delChan <- &triggeredEvent{e: testElement}

	go close(addChan)
	go close(modChan)
	go close(delChan)

	for i := 0; i < 3; i++ {

		if ok := <-notifyChan; !ok {
			t.Error("Expected true")
		}
	}

}
