package filesystem

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type MemoryFileData struct {
	data string
}

// MemoryFile provides an in-memory implementation of a file
type MemoryFile struct {
	content *MemoryFileData
	ch      chan<- bool
}

type MemoryFileEntry struct {
	mu       sync.RWMutex
	file     MemoryFile
	watchers map[chan bool]struct{}
}

// Close the file
func (MemoryFile) Close() error {
	return nil
	// panic("not implemented") // TODO: Implement
}

// Read the file
func (f MemoryFile) Read(p []byte) (n int, err error) {
	n = copy(p, []byte(f.content.data))
	return n, io.EOF
}

func (f MemoryFile) Write(b []byte) (n int, err error) {
	f.content.data = string(b)
	f.ch <- true
	return len(f.content.data), nil
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

func NewMemoryFileSystem() *MemoryFileSystem {
	return &MemoryFileSystem{
		make(map[string]*MemoryFileEntry),
	}
}

// Create a file
func (fs MemoryFileSystem) Create(name string) (File, error) {
	entry, ok := fs.files[name]
	if !ok {
		ch := make(chan bool)

		fs.files[name] = &MemoryFileEntry{
			sync.RWMutex{},
			MemoryFile{
				&MemoryFileData{},
				ch,
			},
			make(map[chan bool]struct{}),
		}
		entry = fs.files[name]

		go func() {
			for range ch {
				for watcher := range entry.watchers {
					watcher <- true
				}
			}
		}()
	}
	return entry.file, nil
}

// Open a file
func (fs MemoryFileSystem) Open(name string) (File, error) {
	entry, ok := fs.files[name]
	if !ok {
		return nil, fmt.Errorf("File not found: %q", name)
	}
	return entry.file, nil
}

// Stat - Not implemented
func (MemoryFileSystem) Stat(name string) (os.FileInfo, error) { panic("not implemented") }

// Watch a file for changes, returns a receiving channel for notifying of change events
func (fs MemoryFileSystem) Watch(name string, done <-chan bool) (<-chan bool, error) {
	files := make(chan bool)

	entry, ok := fs.files[name]
	if !ok {
		return nil, fmt.Errorf("File not found: %q", name)
	}

	entry.mu.Lock()
	entry.watchers[files] = struct{}{}
	entry.mu.Unlock()

	go func() {
		<-done
		go func() {
			for range files {
			}
		}()
		entry.mu.Lock()
		delete(entry.watchers, files)
		entry.mu.Unlock()
		close(files)
	}()

	return files, nil
}
