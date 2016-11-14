package poller

import (
	"os"
	"time"
)

// Element is the interface which describes a Remote Element.
// An Element could be a file, struct or anything which could satisfy this interface.
// Elements are compared between cycles, using cached (previous cycle) and current cycle.
// Elements will be compared if the same key exists in both cached and current cycle.
type Element interface {
	// key for the element
	Name() string

	// modification time
	LastModified() time.Time

	// is a directory
	IsDirectory() bool
}

// A FileElement is an implementation of Element.
// FileElement stores os.FileInfo and can be used as
// for local filesystem files
type FileElement struct {
	os.FileInfo
}

// Name of file returned by FileInfo
func (f *FileElement) Name() string {
	return f.FileInfo.Name()
}

// modified time of file returned by FileInfo
func (f *FileElement) LastModified() time.Time {

	return f.FileInfo.ModTime()
}

// checks if directory returned by FileInfo
func (f *FileElement) IsDirectory() bool {

	return f.FileInfo.IsDir()
}
