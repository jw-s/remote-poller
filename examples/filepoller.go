package main

import (
	"github.com/joelw-s/remote-poller"
	"log"
	"time"
)

type stdOutListener struct{}

func (listener *stdOutListener) OnFileAdded(event poller.Event) {
	log.Printf("%s has been added", event.TriggerCause().Name())
	time.Sleep(2 * time.Minute)
	log.Println("Exiting added event!")
}

func (listener *stdOutListener) OnFileDeleted(event poller.Event) {
	log.Printf("%s has been deleted", event.TriggerCause().Name())
	time.Sleep(3 * time.Minute)
	log.Println("Exiting deleted event!")
}

func (listener *stdOutListener) OnFileModified(event poller.Event) {
	log.Printf("%s has been modified", event.TriggerCause().Name())
	time.Sleep(5 * time.Minute)
	log.Println("Exiting modified event!")
}
func main() {

	//Creates new Polled filesystem directory, which is recursive
	polledDirectory, err := poller.NewFileDirectory("/test")

	if err != nil {
		panic(err)
	}

	filesystemPoller := poller.NewPoller(time.Duration(15*time.Second), polledDirectory, []poller.Receiver{&stdOutListener{}})

	//Starts poller non blocking
	filesystemPoller.Start()

	time.Sleep(2 * time.Minute)

	//Stop shutsdown poller + blocks until all event processing has completed
	filesystemPoller.Stop()
}
