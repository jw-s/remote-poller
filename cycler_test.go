package poller

import (
	"testing"
	"time"
)

type testPolledDirectory struct {
	elements []Element
}

func (pd *testPolledDirectory) ListFiles() ([]Element, error) {
	return pd.elements, nil

}

func createTime(t string) time.Time {
	tm, _ := time.Parse(time.Kitchen, t)
	return tm
}

func TestPollCycle_NotifyFirstCycle(t *testing.T) {

	elements := []Element{
		&testElement{name: "1"},
	}

	pd := testPolledDirectory{elements}

	pc := pollCycle{
		firstRun:        true,
		polledDirectory: &pd,
		cachedElements:  make(chan map[string]Element, 1),
	}

	//trigger first run, gets initial cache
	pc.Notify()

	cached := <-pc.cachedElements

	for _, e := range elements {
		if _, ok := cached[cachedKeyNameFormat(e)]; !ok {
			t.Errorf("%s should exist in cache", e.Name())
		}
	}

}

func TestPollCycle_NotifyDeleted(t *testing.T) {

	elements := []Element{
		&testElement{name: "1"},
		&testElement{name: "2"},
	}

	pd := testPolledDirectory{elements}

	pc := pollCycle{
		firstRun:        true,
		polledDirectory: &pd,
		cachedElements:  make(chan map[string]Element, 1),
		em: &eventTriggerManager{
			receivers: []Receiver{
				testReceiver{},
			}},
	}

	//trigger first run, gets initial cache
	pc.Notify()

	pd.elements = append(elements[:len(elements)-1])

	// trigger another run
	pc.Notify()

	cached := <-pc.cachedElements

	if e, ok := cached["2"]; ok {
		t.Errorf("%s shouldn't exist in cache", e.Name())
	}

}

func TestPollCycle_NotifyAdded(t *testing.T) {

	elements := []Element{
		&testElement{name: "1"},
		&testElement{name: "2"},
	}

	pd := testPolledDirectory{elements}

	pc := pollCycle{
		firstRun:        true,
		polledDirectory: &pd,
		cachedElements:  make(chan map[string]Element, 1),
		em: &eventTriggerManager{
			receivers: []Receiver{
				testReceiver{},
			}},
	}

	//trigger first run, gets initial cache
	pc.Notify()

	toBeAddedElement := &testElement{name: "3"}
	pd.elements = append(elements, toBeAddedElement)
	// trigger another run
	pc.Notify()

	cached := <-pc.cachedElements

	if _, ok := cached["3_false"]; !ok {
		t.Errorf("Element name %s should have been added and exist in cache", toBeAddedElement.Name())
	}

}

func TestPollCycle_NotifyModified(t *testing.T) {

	elements := []Element{
		&testElement{name: "1"},
		&testElement{name: "2"},
	}

	pd := testPolledDirectory{elements}

	pc := pollCycle{firstRun: true,
		polledDirectory: &pd,
		cachedElements:  make(chan map[string]Element, 1),
		em: &eventTriggerManager{receivers: []Receiver{
			testReceiver{},
		}},
	}

	//trigger first run, gets initial cache
	pc.Notify()

	elements[0] = &testElement{name: "1", lastModified: createTime("12:00PM")}

	// trigger another run
	pc.Notify()

	cached := <-pc.cachedElements
	if e, ok := cached["1"]; ok {
		if "0000-01-01 12:00:00 +0000 UTC" != e.LastModified().String() {
			t.Errorf("%d", e.LastModified())
		}

	}

}
