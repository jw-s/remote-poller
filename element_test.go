package remote_poller

import (
	"testing"
	"time"
)

type testElement struct {
	name         string
	lastModified time.Time
	isDirectory  bool
}

func (te *testElement) Name() string {
	return te.name
}
func (te *testElement) LastModified() time.Time {
	return te.lastModified
}
func (te *testElement) IsDirectory() bool {
	return te.isDirectory
}

func TestElementName(t *testing.T) {
	e := testElement{name: "e"}
	if actual := e.name; "e" != actual {
		t.Errorf("Expected e and got %s", actual)
	}
}

func TestElementLastModified(t *testing.T) {
	var date time.Time

	e := testElement{lastModified: date}

	if actual := e.lastModified; date != actual {
		t.Errorf("Expected 0001-01-01 00:00:00 +0000 UTC and got %s", actual)
	}
}

func TestElementIsDirectory(t *testing.T) {

	e := testElement{isDirectory: false}

	if actual := e.isDirectory; false != actual {
		t.Errorf("Expected false and got %s", actual)
	}
}
