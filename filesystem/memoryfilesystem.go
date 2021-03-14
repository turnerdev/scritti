package filesystem

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

// MemoryFile provides an in-memory implementation of a file
type MemoryFile struct {
	entry  *MemoryFileEntry
	buffer *bytes.Buffer
}

type MemoryFileEntry struct {
	mu       sync.RWMutex
	content  string
	watchers map[chan bool]struct{}
}

// Close the file
func (p MemoryFile) Close() error {
	return nil
	// panic("not implemented") // TODO: Implement
}

// Read the file
func (f MemoryFile) Read(p []byte) (n int, err error) {
	return f.buffer.Read(p)
}

func (f MemoryFile) Write(b []byte) (n int, err error) {
	f.entry.mu.Lock()
	defer f.entry.mu.Unlock()
	f.buffer.Reset()
	n, err = f.buffer.Write(b)
	f.entry.content = string(b)
	for watcher := range f.entry.watchers {
		watcher <- true
	}
	return n, err
}

// ReadAt a given offset with a file
func (f MemoryFile) ReadAt(p []byte, off int64) (n int, err error) {
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
		// ch := make(chan bool)

		fs.files[name] = &MemoryFileEntry{
			sync.RWMutex{},
			"",
			make(map[chan bool]struct{}),
		}
		entry = fs.files[name]
	}
	return MemoryFile{
		entry,
		bytes.NewBufferString(entry.content),
	}, nil
}

// Open a file
func (fs MemoryFileSystem) Open(name string) (File, error) {
	entry, ok := fs.files[name]
	if !ok {
		return nil, fmt.Errorf("File not found: %q", name)
	}
	return MemoryFile{
		entry,
		bytes.NewBufferString(entry.content),
	}, nil
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
