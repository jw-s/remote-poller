package remote_poller

// Receiver can be registered as a listener to a PolledDirectory.
//
// Receivers get notified by an eventManager based on what even has occured.
//
// There can be three types of events:
//
// Event for files added, eventManager will call onFileAdded for all listeners
//
// Event for files deleted, eventManager will call onFileDeleted for all listeners
//
// Event for files modified, eventManager will call onFileModified for all listeners
//
// It is up to the implementation what happens once these events are triggered.
type Receiver interface {
	// Processes file added event
	OnFileAdded(Event)

	// Processes file deleted event
	OnFileDeleted(Event)

	// Processes file modified event
	OnFileModified(Event)
}
