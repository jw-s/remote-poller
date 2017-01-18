package poller

// Filter provides a way of filtering Elements from PolledDirectory cycle.
// An Element matching the filter is accepted and carries on as normal.
// If it doesn't match then their lifecycle terminates.
type Filter interface {
	//Accept the Element if true is returned
	Accept(Element) bool
}

type defaultFilter struct{}

func (defaultFilter defaultFilter) Accept(element Element) bool {
	return true
}
