package poller

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

	var listeners []Receiver

	notifyChan := make(chan bool)

	r := testReceiver{notify: notifyChan}

	listeners = append(listeners, r, r)

	em := eventTriggerManager{receivers: listeners}

	testEvent := &triggeredEvent{
		&testElement{name: "testElement"},
	}

	go em.OnFileAdded(testEvent)
	go em.OnFileModified(testEvent)
	go em.OnFileDeleted(testEvent)

	for i := 0; i < 3; i++ {

		if ok := <-notifyChan; !ok {
			t.Error("Expected true")
		}
	}

}
