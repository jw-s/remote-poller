package remote_poller

import (
	"log"
	"sync"
)

// Cycler is an interface used to notify channels of changes to Elements.
type cycler interface {
	Notify(chan<- Event, chan<- Event, chan<- Event) error
}

// PolledDirectory is the interface used to provide a listing of Files.
// ListFiles must return both []Elements and possibly an error.
type PolledDirectory interface {
	ListFiles() ([]Element, error)
}

type pollCycle struct {
	firstRun        bool
	polledDirectory PolledDirectory
	cachedElements  chan map[string]Element
}

func (pc *pollCycle) detectDeletedFiles(del chan<- Event, mod chan<- Event, cachedElements, newElements map[string]Element) {
	defer handleClientError()

	for k, v := range cachedElements {

		if ne, ok := newElements[k]; !ok {
			del <- &triggeredEvent{v}
		} else if !v.LastModified().Equal(ne.LastModified()) {
			mod <- &triggeredEvent{v}
		}
	}
}

func (pc *pollCycle) detectAddedFiles(add chan<- Event, cachedElements, newElements map[string]Element) {

	defer handleClientError()

	for k, v := range newElements {
		if _, ok := cachedElements[k]; !ok {
			add <- &triggeredEvent{v}
		}
	}
}

func (pc *pollCycle) Notify(add chan<- Event, mod chan<- Event, del chan<- Event) error {
	var wg sync.WaitGroup

	if pc.firstRun {

		initialElements := make(map[string]Element)
		elements, err := pc.polledDirectory.ListFiles()

		if err != nil {
			return err
		}

		for _, e := range elements {
			initialElements[e.Name()] = e
		}

		pc.firstRun = false

		pc.cachedElements <- initialElements

		return nil
	}

	cachedElements := <-pc.cachedElements

	newElements := make(map[string]Element)

	listedFiles, err := pc.polledDirectory.ListFiles()

	if err != nil {
		log.Printf("Client implementation has thrown error when listing files : %s", err.Error())
	}

	for _, e := range listedFiles {
		newElements[e.Name()] = e
	}

	wg.Add(2)
	go func(cachedElements, newElements map[string]Element) {
		defer wg.Done()
		pc.detectDeletedFiles(del, mod, cachedElements, newElements)
	}(cachedElements, newElements)

	go func(cachedElements, newElements map[string]Element) {
		defer wg.Done()
		pc.detectAddedFiles(add, cachedElements, newElements)
	}(cachedElements, newElements)

	wg.Wait()

	pc.cachedElements <- newElements

	return nil
}

func handleClientError() {
	if err := recover(); err != nil {
		log.Printf("client has panicked: %s", err)
	}
}
