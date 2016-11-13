package remote_poller

import (
	"log"
	"sync"
)

type Cycler interface {
	Notify(chan<- Event, chan<- Event, chan<- Event) error
}

type PolledDirectory interface {
	ListFiles() ([]Element, error)
}

type pollCycle struct {
	firstRun        bool
	mux             sync.Mutex
	polledDirectory PolledDirectory
	cachedElements  map[string]Element
}

func (pc *pollCycle) detectDeletedFiles(del chan<- Event, mod chan<- Event, newElements map[string]Element) {
	defer handleClientError()

	for k, v := range pc.cachedElements {
		if ne, ok := newElements[k]; !ok {
			del <- &triggeredEvent{v}
		} else {
			if !v.LastModified().Equal(ne.LastModified()) {
				mod <- &triggeredEvent{v}
			}
		}
	}
}

func (pc *pollCycle) detectAddedFiles(add chan<- Event, newElements map[string]Element) {

	defer handleClientError()

	for k, v := range newElements {
		if _, ok := pc.cachedElements[k]; !ok {
			add <- &triggeredEvent{v}
		}
	}
}

func (pc *pollCycle) Notify(add chan<- Event, del chan<- Event, mod chan<- Event) error {
	var wg sync.WaitGroup

	defer pc.mux.Unlock()
	pc.mux.Lock()

	if pc.firstRun {

		pc.cachedElements = make(map[string]Element)
		elements, err := pc.polledDirectory.ListFiles()

		if err != nil {
			return err
		}

		for _, e := range elements {
			pc.cachedElements[e.Name()] = e
		}

		pc.firstRun = false

		return nil
	}

	newElements := make(map[string]Element)

	listedFiles, err := pc.polledDirectory.ListFiles()

	if err != nil {
		log.Printf("Client implementation has thrown error when listing files : %s", err.Error())
	}

	for _, e := range listedFiles {
		newElements[e.Name()] = e
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		pc.detectDeletedFiles(del, mod, newElements)
	}()

	go func() {
		defer wg.Done()
		pc.detectAddedFiles(add, newElements)
	}()

	wg.Wait()

	pc.cachedElements = newElements

	return nil
}

func handleClientError() {
	if err := recover(); err != nil {
		log.Printf("client has panicked: %s", err)
	}
}
