package remote_poller

import (
	"testing"
	"time"
)

type testTicker struct {
	ticker chan time.Time
}

func (t *testTicker) Tick() <-chan time.Time { return t.ticker }

func (t *testTicker) Stop() { close(t.ticker) }

func TestPoller_Start(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			t.Error("Channel should be open and listening")
		}
	}()

	add, del, mod, tickerChan := make(chan Event), make(chan Event), make(chan Event), make(chan time.Time)

	tc := &triggerChannels{add: add, del: del, mod: mod}
	ticker := &testTicker{tickerChan}
	elements := []Element{&testElement{name: "1"}, &testElement{name: "2"}}
	pd := &testPolledDirectory{elements}
	pc := &pollCycle{polledDirectory: pd, cachedElements: make(chan map[string]Element, 1)}

	notifyChan := make(chan bool)

	listeners := []Receiver{testReceiver{notify: notifyChan}}

	em := EventTriggerManager{receivers: listeners}

	poller := poller{tc: tc, ticker: ticker, cycler: pc, em: &em}

	poller.Start()

	ticker.ticker <- time.Now()

}

func TestPoller_Stop(t *testing.T) {

	defer func() {
		if err := recover(); err == nil {
			t.Error("Channel should be closed")
		}
	}()
	add, del, mod, tickerChan := make(chan Event), make(chan Event), make(chan Event), make(chan time.Time)

	tc := &triggerChannels{add: add, del: del, mod: mod}
	ticker := &testTicker{tickerChan}
	elements := []Element{&testElement{name: "1"}, &testElement{name: "2"}}
	pd := &testPolledDirectory{elements}
	pc := &pollCycle{polledDirectory: pd, cachedElements: make(chan map[string]Element, 1)}

	notifyChan := make(chan bool)

	listeners := []Receiver{testReceiver{notify: notifyChan}}

	em := EventTriggerManager{receivers: listeners}

	poller := poller{tc: tc, ticker: ticker, cycler: pc, em: &em}

	poller.Start()
	poller.Stop()
	ticker.ticker <- time.Now()

}
