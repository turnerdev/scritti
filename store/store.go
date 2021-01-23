package store

import (
	"io"
	"io/ioutil"
	"os"
	"scritti/core"
)

var fs fileSystem = osFS{}

type fileSystem interface {
	Open(name string) (file, error)
	Stat(name string) (os.FileInfo, error)
}

type file interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (file, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

// IComponentStore provides an interface to retrieve components
type IComponentStore interface {
	Get(name string) (*core.Component, error)
}

// FileStore TODO
type FileStore struct {
	// Reader io.Reader,
	fs    fileSystem
	cache map[string]string
}

// NewFileStore returns a new File Store
func NewFileStore() *FileStore {
	return &FileStore{
		osFS{},
		map[string]string{},
	}
}

// Get TODO
func (c *FileStore) Get(name string) (*core.Component, error) {
	var component *core.Component
	if val, ok := c.cache[name]; ok {
		component = core.ParseComponent(val)
	} else {
		// var err error
		file, err := c.fs.Open("sampledata/" + name)
		if err != nil {
			return component, nil
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return component, nil
		}

		source := string(data)
		c.cache[name] = source
		component = core.ParseComponent(source)
	}
	return component, nil
}

// Expand a component
func Expand(component *core.Component, store IComponentStore) *core.Component {
	return component
}

// import (
// 	"io/ioutil"
// )

// type IRepository interface {
// 	get() <- chan (string, err)
// }
//

// type fileRepository struct {
// }

// func loadComponent(name string) {

// }

//make([]Element, 10)
// if err != nil {
//     fmt.Println("File reading error", err)
//     return
// }
// fmt.Println("Contents of file:", string(data))

// ReadFile blah
