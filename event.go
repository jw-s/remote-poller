package poller

import "sync"

// Event interface used to provide listeners with triggerCause.
type Event interface {
	TriggerCause() Element // returns the Element which trigger Event
}

type triggeredEvent struct {
	e Element
}

func (te triggeredEvent) TriggerCause() Element {
	return te.e
}

// Event manager interface is used to handle and forward events passed from higher up,
// to all registered listeners.
type eventManager interface {
	// Receives event and processes Files which have been added
	OnFileAdded(Event)

	// Receives event and processes Files which have been deleted
	OnFileDeleted(Event)

	// Receives event and processes Files which have been modified
	OnFileModified(Event)

	// Waits for all events to be processed before shutting down
	ShutDownAndWait()
}

type eventTriggerManager struct {
	receivers []Receiver
	wg        sync.WaitGroup
}

func (em *eventTriggerManager) OnFileAdded(event Event) {

	for _, r := range em.receivers {
		em.wg.Add(1)
		go func() {
			defer em.wg.Done()
			r.OnFileAdded(event)
		}()
	}

}

func (em *eventTriggerManager) OnFileDeleted(event Event) {

	for _, r := range em.receivers {
		em.wg.Add(1)
		go func() {
			defer em.wg.Done()
			r.OnFileDeleted(event)
		}()
	}

}

func (em *eventTriggerManager) OnFileModified(event Event) {

	for _, r := range em.receivers {
		em.wg.Add(1)
		go func() {
			defer em.wg.Done()
			r.OnFileModified(event)
		}()
	}
}

func (em *eventTriggerManager) ShutDownAndWait() {
	em.wg.Wait()
}
