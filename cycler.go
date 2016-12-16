package poller

import (
	"fmt"
	"log"
	"sync"
)

// Cycler is an interface used to notify channels of changes to Elements.
type cycler interface {
	Notify() error
	Stop()
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

func (pc *pollCycle) detectDeletedFilesAndNotify(cachedElements, newElements map[string]Element) {
	defer handleClientError()

	for k, v := range cachedElements {

		if ne, ok := newElements[k]; !ok {

			go pc.em.OnFileDeleted(triggeredEvent{v})

		} else if !v.LastModified().Equal(ne.LastModified()) {

			go pc.em.OnFileModified(&triggeredEvent{v})

		}
	}
}

func (pc *pollCycle) detectAddedFilesAndNotify(cachedElements, newElements map[string]Element) {

	defer handleClientError()

	for k, v := range newElements {

		if _, ok := cachedElements[k]; !ok {

			go pc.em.OnFileAdded(&triggeredEvent{v})

		}
	}
}

func (pc *pollCycle) onFirstRun() error {

	initialElements := make(map[string]Element)
	elements, err := pc.polledDirectory.ListFiles()

	if err != nil {
		return err
	}

	for _, e := range elements {
		initialElements[cachedKeyNameFormat(e)] = e
	}

	pc.firstRun = false

	pc.cachedElements <- initialElements

	return nil

}

func (pc *pollCycle) Notify() error {
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

		newElements[cachedKeyNameFormat(e)] = e
	}

	wg.Add(2)
	go func(cachedElements, newElements map[string]Element) {
		defer wg.Done()

		pc.detectDeletedFilesAndNotify(cachedElements, newElements)

	}(cachedElements, newElements)

	go func(cachedElements, newElements map[string]Element) {
		defer wg.Done()

		pc.detectAddedFilesAndNotify(cachedElements, newElements)

	}(cachedElements, newElements)

	wg.Wait()

	pc.cachedElements <- newElements

	return nil
}

func (pc *pollCycle) Stop() {
	pc.em.ShutDownAndWait()
}

func handleClientError() {
	if err := recover(); err != nil {
		log.Printf("client has panicked: %s", err)
	}
}

func cachedKeyNameFormat(element Element) string {
	return fmt.Sprintf("%s_%t", element.Name(), element.IsDirectory())

}
