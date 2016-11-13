package remote_poller

type Receiver interface {
	OnFileAdded(Event)
	OnFileDeleted(Event)
	OnFileModified(Event)
}
