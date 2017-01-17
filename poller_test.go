package poller

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

	tickerChan := make(chan time.Time)

	ticker := &testTicker{tickerChan}

	elements := []Element{
		&testElement{name: "1"},
		&testElement{name: "2"},
	}

	pd := &testPolledDirectory{elements}

	notifyChan := make(chan bool)

	listeners := []Receiver{
		testReceiver{
			notify: notifyChan,
		},
	}

	pc := &pollCycle{
		polledDirectory: pd,
		cachedElements:  make(chan map[string]Element, 1),
		em: &eventTriggerManager{
			receivers: listeners,
		},
	}

	poller := poller{
		ticker: ticker,
		cycler: pc,
	}

	poller.Start()

	ticker.ticker <- time.Now()

}

func TestPoller_Stop(t *testing.T) {

	defer func() {
		if err := recover(); err == nil {
			t.Error("Channel should be closed")
		}
	}()
	tickerChan := make(chan time.Time)

	ticker := &testTicker{
		ticker: tickerChan,
	}

	elements := []Element{
		&testElement{name: "1"},
		&testElement{name: "2"},
	}

	pd := &testPolledDirectory{elements}

	notifyChan := make(chan bool)

	listeners := []Receiver{
		testReceiver{
			notify: notifyChan,
		},
	}

	pc := &pollCycle{
		polledDirectory: pd,
		cachedElements:  make(chan map[string]Element, 1),
		em: &eventTriggerManager{
			receivers: listeners,
		},
	}

	poller := poller{
		ticker: ticker,
		cycler: pc,
	}

	poller.Start()
	poller.Stop()
	ticker.ticker <- time.Now()

}
