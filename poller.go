package poller

import (
	"log"
	"time"
)

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
	ticker ticker
	cycler cycler
}

// Creates a poller used to trigger the Cycler at specified interval.
func NewPoller(d time.Duration, pollDir PolledDirectory, listeners []Receiver) *poller {

	return NewPollerWithFilters(d, pollDir, listeners, nil)

}

// Creates a poller used to trigger the Cycler at specified interval.
func NewPollerWithFilters(d time.Duration, pollDir PolledDirectory, listeners []Receiver, filters []Filter) *poller {

	if len(filters) == 0 && cap(filters) == 0 {

		filters = []Filter{
			defaultFilter{},
		}
	}

	cycler := pollCycle{
		firstRun:        true,
		polledDirectory: pollDir,
		cachedElements:  make(chan map[string]Element, 1),
		em: &eventTriggerManager{
			receivers: listeners,
			filters:   filters,
		},
	}

	return &poller{ticker: &pollTicker{time.NewTicker(d)},
		cycler: &cycler}
}

// Starts the poller and triggers cycle at set interval.
func (p *poller) Start() {

	log.Println("Starting poller")
	log.Println("Will start polling after initial tick...")

	ticker := p.ticker
	go func() {
		for {

			_, open := <-ticker.Tick()

			log.Println("Ticking...")

			if !open {
				return
			}

			go func() {
				err := p.cycler.Notify()
				if err != nil {
					log.Fatalf("Client has thrown error, exiting... %s", err.Error())
				}
			}()

		}
	}()
}

// Stops the poller.
func (p *poller) Stop() {
	log.Println("Stopping poller")
	p.ticker.Stop()
	p.cycler.Stop()
}
