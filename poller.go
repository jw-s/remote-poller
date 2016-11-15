package poller

import (
	"log"
	"os"
	"time"
)

type triggerChannels struct {
	add, mod, del chan Event
}

type ticker interface {
	Tick() <-chan time.Time
	Stop()
}

type pollTicker struct {
	*time.Ticker
}

func (t *pollTicker) Tick() <-chan time.Time { return t.C }

func (t *pollTicker) Stop() { t.Ticker.Stop() }

type poller struct {
	tc     *triggerChannels
	ticker ticker
	cycler cycler
}

// Creates a poller used to trigger the Cycler at specified interval.
func NewPoller(d time.Duration, pollDir PolledDirectory, listeners []Receiver) *poller {

	tc := triggerChannels{make(chan Event, 1), make(chan Event, 1), make(chan Event, 1)}
	cycler := pollCycle{firstRun: true,
		polledDirectory: pollDir,
		cachedElements:  make(chan map[string]Element, 1),
		em:              &eventTriggerManager{listeners}}

	return &poller{&tc, &pollTicker{time.NewTicker(d)}, &cycler}

}

// Starts the poller and triggers cycle at set interval.
func (p *poller) Start() {
	log.SetOutput(os.Stdout)
	log.Println("Starting poller")
	log.Println("Will start polling after initial tick...")
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

// Stops the poller.
func (p *poller) Stop() {
	p.ticker.Stop()
	close(p.tc.mod)
	close(p.tc.add)
	close(p.tc.del)
}
