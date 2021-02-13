package filesystem

import (
	"errors"
	"io"
	"os"
	"sync"
)

// MemoryFile provides an in-memory implementation of a file
type MemoryFile struct {
	content string
}

type MemoryFileEntry struct {
	mu       sync.RWMutex
	file     MemoryFile
	watchers map[chan File]struct{}
}

// Close the file
func (MemoryFile) Close() error {
	panic("not implemented") // TODO: Implement
}

// Read the file
func (f MemoryFile) Read(p []byte) (n int, err error) {
	n = copy(p, []byte(f.content))
	return n, io.EOF
}

// ReadAt a given offset with a file
func (MemoryFile) ReadAt(p []byte, off int64) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

// Seek the file stub to a given position
func (MemoryFile) Seek(offset int64, whence int) (int64, error) {
	panic("not implemented") // TODO: Implement
}

// Stat the file stub
func (MemoryFile) Stat() (os.FileInfo, error) {
	panic("not implemented") // TODO: Implement
}

// MemoryFileSystem implements an in-memory File System
type MemoryFileSystem struct {
	files map[string]*MemoryFileEntry
}

func (fs MemoryFileSystem) Write(name string, data string) {
	entry, ok := fs.files[name]
	if !ok {
		fs.files[name] = &MemoryFileEntry{
			sync.RWMutex{},
			MemoryFile{data},
			make(map[chan File]struct{}),
		}
	}
	entry.mu.Lock()
	defer entry.mu.Unlock()
	entry.file.content = data
}
func (fs MemoryFileSystem) Open(name string) (File, error) {
	entry, ok := fs.files[name]
	if !ok {
		return nil, errors.New("File not found")
	}
	return entry.file, nil
}
func (MemoryFileSystem) Stat(name string) (os.FileInfo, error) { panic("not implemented") }
func (fs MemoryFileSystem) Watch(name string, done <-chan bool) <-chan File {
	files := make(chan File)

	entry, ok := fs.files[name]
	if !ok {
		panic("File not found")
	}

	entry.watchers[files] = struct{}{}

	go func() {
		for {
			select {
			case <-done:
				delete(entry.watchers, files)
				close(files)
				return
			}
		}
	}()

	return files
}
