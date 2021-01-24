package core

import (
	"io"
	"os"
)

// StubFS provides a stubbed file system
type StubFS struct {
	files map[string]StubFile
}

// StubFile provides a stubbed file
type StubFile struct {
	content string
}

// Close the file stub
func (StubFile) Close() error {
	panic("not implemented") // TODO: Implement
}

// Read the file stub
func (f StubFile) Read(p []byte) (n int, err error) {
	n = copy(p, []byte(f.content))
	return n, io.EOF
}

// ReadAt a given offset with a file stub
func (StubFile) ReadAt(p []byte, off int64) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

// Seek the file stub to a given position
func (StubFile) Seek(offset int64, whence int) (int64, error) {
	panic("not implemented") // TODO: Implement
}

// Stat the file stub
func (StubFile) Stat() (os.FileInfo, error) {
	panic("not implemented") // TODO: Implement
}

// Open a file with the stubbed file system
func (fs StubFS) Open(name string) (File, error) { return fs.files[name], nil }

// Stat a file in the stubbed file system
func (StubFS) Stat(name string) (os.FileInfo, error) {
	panic("not implemented") // TODO: Implement
}
