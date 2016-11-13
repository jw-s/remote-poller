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
	*os.File
}

func (f *FileElement) Name() string {
	return f.File.Name()
}

func (f *FileElement) LastModified() time.Time {
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	return info.ModTime()
}

func (f *FileElement) IsDirectory() bool {
	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	return info.IsDir()
}
