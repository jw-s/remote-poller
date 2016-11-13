package remote_poller

type Event interface {
	TriggerCause() Element
}

type triggeredEvent struct {
	e Element
}

func (te triggeredEvent) TriggerCause() Element {
	return te.e
}

type EventManager interface {
	OnFileAdded(<-chan Event)
	OnFileDeleted(<-chan Event)
	OnFileModified(<-chan Event)
}

type EventTriggerManager struct {
	receivers []Receiver
}

func (em *EventTriggerManager) OnFileAdded(eventChan <-chan Event) {
	for {
		select {
		case event := <-eventChan:
			for _, r := range em.receivers {
				go r.OnFileAdded(event)
			}
		}

	}
}

func (em *EventTriggerManager) OnFileDeleted(eventChan <-chan Event) {
	for {
		select {
		case event := <-eventChan:
			for _, r := range em.receivers {
				go r.OnFileDeleted(event)
			}
		}

	}
}

func (em *EventTriggerManager) OnFileModified(eventChan <-chan Event) {
	for {
		select {
		case event := <-eventChan:
			for _, r := range em.receivers {
				go r.OnFileModified(event)
			}
		}

	}
}
