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
	em     eventManager
}

// Creates a poller used to trigger the Cycler at specified interval.
func NewPoller(d time.Duration, pollDir PolledDirectory, listeners []Receiver) *poller {

	tc := &triggerChannels{make(chan Event), make(chan Event), make(chan Event)}
	cycler := pollCycle{firstRun: true, polledDirectory: pollDir, cachedElements: make(chan map[string]Element)}

	return &poller{tc, &pollTicker{time.NewTicker(d)}, &cycler, &eventTriggerManager{listeners}}

}

// Starts the poller and triggers cycle at set interval.
func (p *poller) Start() {
	log.SetOutput(os.Stdout)
	log.Println("Starting poller")
	log.Println("Will start polling after initial tick...")
	add, mod, del := p.tc.add, p.tc.mod, p.tc.del

	go p.em.OnFileAdded(add)
	go p.em.OnFileModified(mod)
	go p.em.OnFileDeleted(del)

	ticker := p.ticker
	go func() {
		for {
			select {
			case _, open := <-ticker.Tick():
				log.Println("Starting poll cycle")
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
}
