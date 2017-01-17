package poller

import (
	"strings"
	"testing"
)

type testFilter struct{}

func (testFilter testFilter) Accept(e Element) bool {
	return strings.Contains(e.Name(), "test")
}

func TestFilter_Filter(t *testing.T) {

	filterNameTests := []struct {
		name   string
		accept bool
	}{
		{"testElement", true},
		{"Elementtest", true},
		{"Element", false},
		{"tset", false},
	}

	tf := testFilter{}

	for _, filterTest := range filterNameTests {

		te := testElement{
			name: filterTest.name,
		}
		if expected, actual := filterTest.accept, tf.Accept(&te); expected != actual {

			t.Errorf("Expected %t, got %t", expected, actual)
		}

	}

}
