package remote_poller

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
	OnFileAdded(<-chan Event)

	// Receives event and processes Files which have been deleted
	OnFileDeleted(<-chan Event)

	// Receives event and processes Files which have been modified
	OnFileModified(<-chan Event)
}

type eventTriggerManager struct {
	receivers []Receiver
}

func (em *eventTriggerManager) OnFileAdded(eventChan <-chan Event) {
	for {
		select {
		case event := <-eventChan:
			for _, r := range em.receivers {
				go r.OnFileAdded(event)
			}
		}

	}
}

func (em *eventTriggerManager) OnFileDeleted(eventChan <-chan Event) {
	for {
		select {
		case event := <-eventChan:
			for _, r := range em.receivers {
				go r.OnFileDeleted(event)
			}
		}

	}
}

func (em *eventTriggerManager) OnFileModified(eventChan <-chan Event) {
	for {
		select {
		case event := <-eventChan:
			for _, r := range em.receivers {
				go r.OnFileModified(event)
			}
		}

	}
}
