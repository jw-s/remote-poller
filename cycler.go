package poller

import (
	"log"
	"sync"
)

// Cycler is an interface used to notify channels of changes to Elements.
type cycler interface {
	Notify(chan Event, chan Event, chan Event) error
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
	em              eventManager
}

func (pc *pollCycle) detectDeletedFilesAndNotify(del chan Event, mod chan Event, cachedElements, newElements map[string]Element) {
	defer handleClientError()

	for k, v := range cachedElements {

		if ne, ok := newElements[k]; !ok {

			go sendEventToChannel(del, v)
			go pc.em.OnFileDeleted(del)

		} else if !v.LastModified().Equal(ne.LastModified()) {

			go sendEventToChannel(mod, v)
			go pc.em.OnFileModified(mod)

		}
	}
}

func (pc *pollCycle) detectAddedFilesAndNotify(add chan Event, cachedElements, newElements map[string]Element) {

	defer handleClientError()

	for k, v := range newElements {

		if _, ok := cachedElements[k]; !ok {

			go sendEventToChannel(add, v)
			go pc.em.OnFileAdded(add)

		}
	}
}

func sendEventToChannel(eventChan chan<- Event, element Element) {
	eventChan <- &triggeredEvent{element}
}

func (pc *pollCycle) onFirstRun() error {

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

func (pc *pollCycle) Notify(add chan Event, mod chan Event, del chan Event) error {
	var wg sync.WaitGroup

	if pc.firstRun {

		return pc.onFirstRun()
	}

	cachedElements := <-pc.cachedElements

	listedFiles, err := pc.polledDirectory.ListFiles()

	if err != nil {
		log.Printf("Client implementation has thrown error when listing files : %s", err.Error())
		return nil
	}

	newElements := make(map[string]Element)

	for _, e := range listedFiles {
		newElements[e.Name()] = e
	}

	wg.Add(2)
	go func(del, mod chan Event, cachedElements, newElements map[string]Element) {
		defer wg.Done()

		pc.detectDeletedFilesAndNotify(del, mod, cachedElements, newElements)

	}(del, mod, cachedElements, newElements)

	go func(add chan Event, cachedElements, newElements map[string]Element) {
		defer wg.Done()

		pc.detectAddedFilesAndNotify(add, cachedElements, newElements)

	}(add, cachedElements, newElements)

	wg.Wait()

	pc.cachedElements <- newElements

	return nil
}

func handleClientError() {
	if err := recover(); err != nil {
		log.Printf("client has panicked: %s", err)
	}
}
