package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var fs FileSystem = osFS{}

// FileSystem interface
type FileSystem interface {
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
}

// File interface
type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

// osFS implements fileSystem using the local disk.
type osFS struct{}

func (osFS) Open(name string) (File, error)        { return os.Open(name) }
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

// AssetType
type AssetType int

// AssetType enum
const (
	ComponentType AssetType = iota
	StyleType
)

// AssetStore provides an interface to retrieve components
type AssetStore interface {
	Get(assetType AssetType, name string) (Asset, error)
}

// FileStore TODO
type FileStore struct {
	// Reader io.Reader,
	path  string
	fs    FileSystem
	cache map[AssetType]map[string]string
}

// NewFileStore returns a new File Store
func NewFileStore(path string) *FileStore {
	return &FileStore{
		path,
		osFS{},
		map[AssetType]map[string]string{
			ComponentType: {},
			StyleType:     {},
		},
	}
}

var assetPath = map[AssetType]string{
	StyleType: "style",
}

// Get TODO
func (c FileStore) Get(assetType AssetType, name string) (Asset, error) {
	var asset Asset

	// Load asset from file system if not found in cache
	if _, ok := c.cache[assetType][name]; !ok {
		fmt.Printf("Path %q\n", filepath.Join(c.path, assetPath[assetType], name))
		file, err := c.fs.Open(filepath.Join(c.path, assetPath[assetType], name))
		if err != nil {
			return asset, nil
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return asset, nil
		}

		source := string(data)
		c.cache[assetType][name] = source
	}

	switch assetType {
	case ComponentType:
		asset = ParseComponent(c.cache[ComponentType][name])
	case StyleType:
		fmt.Println("PARSE STYLE")
		style := ParseStyle(c.cache[StyleType][name])
		style.name = name
		asset = style
	default:
		panic("not implemented")
	}

	return asset, nil
}

// Expand a component
func Expand(component *Component, store AssetStore) *Component {
	return component
}
