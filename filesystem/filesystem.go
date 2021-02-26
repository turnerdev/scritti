package filesystem

import (
	"io"
	"os"
)

// File interface
type File interface {
	io.Closer
	io.Reader
	io.Writer
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

// FileSystem interface
type FileSystem interface {
	Create(name string) (File, error)
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
	Watch(name string, done <-chan bool) (<-chan bool, error)
}
