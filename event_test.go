package poller

import (
	"sync"
	"testing"
)

type testReceiverChan struct {
	notify      chan<- bool
	delReceived bool
	modReceived bool
}

func (r testReceiverChan) OnFileAdded(e Event) {
	r.notify <- true
}

func (r testReceiverChan) OnFileDeleted(e Event) {
	r.notify <- true
}
func (r testReceiverChan) OnFileModified(e Event) {
	r.notify <- true
}

type testEventReceiver struct {
	results []struct{}
	mut     sync.Mutex
}

func (r *testEventReceiver) addToResultsOnEvent() {
	r.mut.Lock()
	defer r.mut.Unlock()
	r.results = append(r.results, struct{}{})
}

func (r *testEventReceiver) OnFileAdded(e Event) {
	r.addToResultsOnEvent()
}

func (r *testEventReceiver) OnFileDeleted(e Event) {
	r.addToResultsOnEvent()
}
func (r *testEventReceiver) OnFileModified(e Event) {
	r.addToResultsOnEvent()
}

func TestEventTrigger_OnEvents(t *testing.T) {
	triggerTests := []struct {
		name         string
		triggerCount int
	}{
		{"testElement", 3},
		{"Elementtest", 3},
		{"Element", 0},
		{"tset", 0},
	}

	for _, triggerTest := range triggerTests {

		testReceiver := &testEventReceiver{results: nil}

		listeners := []Receiver{
			testReceiver,
		}

		em := eventTriggerManager{
			receivers: listeners,
			filters: []Filter{
				testFilter{},
			},
		}

		testEvent := &triggeredEvent{
			&testElement{
				name: triggerTest.name,
			},
		}

		em.OnFileAdded(testEvent)
		em.OnFileModified(testEvent)
		em.OnFileDeleted(testEvent)

		em.wg.Wait()

		if expected, actual := triggerTest.triggerCount, len(testReceiver.results); expected != actual {

			t.Errorf("Expected %d, got %d", expected, actual)

		}
	}
}
