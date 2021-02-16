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

func NewMemoryFileSystem() *MemoryFileSystem {
	return &MemoryFileSystem{
		make(map[string]*MemoryFileEntry),
	}
}

func (fs MemoryFileSystem) Write(name string, data string) {
	// Fetch entry from filesystem
	entry, ok := fs.files[name]
	if !ok {
		// Create new entry if missing
		fs.files[name] = &MemoryFileEntry{
			sync.RWMutex{},
			MemoryFile{data},
			make(map[chan File]struct{}),
		}
		entry = fs.files[name]
	} else {
		// Update file content
		entry.file.content = data
	}

	// Notify watchers
	entry.mu.RLock()
	for watcher := range fs.files[name].watchers {
		watcher <- fs.files[name].file
	}
	entry.mu.RUnlock()
}

// Open a file
func (fs MemoryFileSystem) Open(name string) (File, error) {
	entry, ok := fs.files[name]
	if !ok {
		return nil, errors.New("File not found")
	}
	return entry.file, nil
}

// Stat - Not implemented
func (MemoryFileSystem) Stat(name string) (os.FileInfo, error) { panic("not implemented") }

// Watch a file for changes, returns a receiving channel for notifying of change events
func (fs MemoryFileSystem) Watch(name string, done <-chan bool) (<-chan File, error) {
	files := make(chan File)

	entry, ok := fs.files[name]
	if !ok {
		return nil, errors.New("File not found")
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
