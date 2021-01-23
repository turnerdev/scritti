package store

import (
	"io"
	"os"
	"testing"
)

func TestFileStore(t *testing.T) {

	t.Run("Test cache miss", func(t *testing.T) {
		componentName := "example"

		fileStore := FileStore{
			stubFS{
				map[string]stubFile{
					"sampledata/test": {componentName},
				},
			},
			map[string]string{},
		}

		component, err := fileStore.Get("test")

		if err != nil {
			t.Errorf("Cache error %q", err)
		}

		if component.GetName() != componentName {
			t.Errorf("Expected component name %q got %q", componentName, component.GetName())
		}

	})

}

// type fidle interface {
// 	io.Closer
// 	io.Reader
// 	io.ReaderAt
// 	io.Seeker
// 	Stat() (os.FileInfo, error)
// }

type stubFS struct {
	files map[string]stubFile
}
type stubFile struct {
	content string
}

func (stubFile) Close() error {
	panic("not implemented") // TODO: Implement
}

func (f stubFile) Read(p []byte) (n int, err error) {
	n = copy(p, []byte(f.content))
	return n, io.EOF
}

func (stubFile) ReadAt(p []byte, off int64) (n int, err error) {
	panic("not implemented") // TODO: Implement
}

func (stubFile) Seek(offset int64, whence int) (int64, error) {
	panic("not implemented") // TODO: Implement
}

func (stubFile) Stat() (os.FileInfo, error) {
	panic("not implemented") // TODO: Implement
}

func (fs stubFS) Open(name string) (file, error) { return fs.files[name], nil }
func (stubFS) Stat(name string) (os.FileInfo, error) {
	panic("not implemented") // TODO: Implement
}
