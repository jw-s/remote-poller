package remote_poller

import (
	"log"
	"time"
)

type triggerChannels struct {
	add chan Event
	mod chan Event
	del chan Event
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

func NewTicker(d time.Duration) *ticker {
	return &ticker{time.NewTicker(d), d}
}

type Poller struct {
	tc     *triggerChannels
	ticker Ticker
	cycler Cycler
}

func (p *Poller) Start() {
	add, mod, del := p.tc.add, p.tc.mod, p.tc.del
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

func (p *Poller) Stop() {
	p.ticker.Stop()
}
