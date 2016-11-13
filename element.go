package remote_poller

import (
	"os"
	"time"
)

type Element interface {
	Name() string
	LastModified() time.Time
	IsDirectory() bool
}

type FileElement struct {
	os.FileInfo
}

func (f *FileElement) Name() string {
	return f.FileInfo.Name()
}

func (f *FileElement) LastModified() time.Time {

	return f.FileInfo.ModTime()
}

func (f *FileElement) IsDirectory() bool {

	return f.FileInfo.IsDir()
}
