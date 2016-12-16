package poller

import (
	"os"
	"path/filepath"
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
	name string
}

// Name of file returned by FileInfo
func (f *FileElement) Name() string {
	return f.name
}

// LastModified time of file returned by FileInfo
func (f *FileElement) LastModified() time.Time {

	return f.FileInfo.ModTime()
}

// IsDirectory returned by FileInfo
func (f *FileElement) IsDirectory() bool {

	return f.FileInfo.IsDir()
}

func (f *FileElement) ListFiles() ([]Element, error) {
	var fileElements []Element
	err := filepath.Walk(f.Name(), func(path string, info os.FileInfo, err error) error {
		fileElements = append(fileElements, &FileElement{FileInfo: info, name: filepath.ToSlash(path)})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileElements, nil
}

//NewFileDirectory creates a PolledDirectory for the specified root dir.
// This implementation handles nested directories.
func NewFileDirectory(rootFilePath string) (PolledDirectory, error) {
	file, err := os.Open(rootFilePath)

	if err != nil {
		return nil, err
	}
	info, err := file.Stat()

	if err != nil {
		return nil, err
	}

	return &FileElement{FileInfo: info, name: rootFilePath}, nil
}
