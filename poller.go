package remote_poller

import (
	"log"
	"time"
)

type triggerChannels struct {
	add, mod, del chan Event
}

type Ticker interface {
	Tick() <-chan time.Time
	Stop()
}

type ticker struct {
	*time.Ticker
	d time.Duration
}

func (t *ticker) Tick() <-chan time.Time { return t.C }

func (t *ticker) Stop() { t.Ticker.Stop() }

func newTicker(d time.Duration) *ticker {
	return &ticker{time.NewTicker(d), d}
}

type poller struct {
	tc     *triggerChannels
	ticker Ticker
	cycler Cycler
	em     EventManager
}

func NewPoller(d time.Duration, pollDir PolledDirectory, listeners []Receiver) *poller {

	tc := &triggerChannels{make(chan Event), make(chan Event), make(chan Event)}
	cycler := pollCycle{firstRun: true, polledDirectory: pollDir}

	return &poller{tc, newTicker(d), &cycler, &EventTriggerManager{listeners}}

}

func (p *poller) Start() {
	add, mod, del := p.tc.add, p.tc.mod, p.tc.del

	go p.em.OnFileAdded(add)
	go p.em.OnFileModified(mod)
	go p.em.OnFileDeleted(del)

	ticker := p.ticker
	go func() {
		for {
			select {
			case _, open := <-ticker.Tick():

				if !open {
					return
				}

				go func() {
					err := p.cycler.Notify(add, mod, del)
					if err != nil {
						log.Fatalf("Client has thrown error, exiting... %s", err.Error())
					}
				}()
			}

		}
	}()
}

func (p *poller) Stop() {
	p.ticker.Stop()
}
