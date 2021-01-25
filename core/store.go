package core

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var fs FileSystem = osFS{}

// FileSystem interface
type FileSystem interface {
	Open(name string) (File, error)
	Stat(name string) (os.FileInfo, error)
	Watch(name string, ch chan<- bool)
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

func (osFS) Open(name string) (File, error)        { return os.Open(name) } // TODO Provide close
func (osFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (osFS) Watch(name string, ch chan<- bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool, 1)

	go func(ch chan<- bool) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				// log.Println("TEST B")
				// ch <- true
				log.Println("TEST A")
				done <- true
				log.Println("EVENT DONE")
			}
		}
	}(ch)

	err = watcher.Add(name)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	log.Println("DONE")
	ch <- true
}

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
	Clear(assetType AssetType, name string)
	Watch(assetType AssetType, name string, ch chan<- bool)
}

// FileStore TODO
type FileStore struct {
	// Reader io.Reader,
	path  string
	fs    FileSystem
	cache map[AssetType]map[string]string
	watch chan<- Asset
}

// NewFileStore returns a new File Store
func NewFileStore(path string) *FileStore {
	return &FileStore{
		path: path,
		fs:   osFS{},
		cache: map[AssetType]map[string]string{
			ComponentType: {},
			StyleType:     {},
		},
	}
}

var assetPath = map[AssetType]string{
	StyleType: "style",
}

func (c FileStore) Clear(assetType AssetType, name string) {
	delete(c.cache[assetType], name)
}

func (c FileStore) Watch(assetType AssetType, name string, ch chan<- bool) {
	filepath := filepath.Join(c.path, assetPath[assetType], name)
	ch2 := make(chan bool)
	log.Printf("Starting to watch %q", filepath)
	c.fs.Watch(filepath, ch2)
	for range ch2 {
		log.Printf("!Watched: %q", filepath)
		c.Clear(assetType, name)
		ch <- true
	}
}

// Get TODO
func (c FileStore) Get(assetType AssetType, name string) (Asset, error) {
	var asset Asset

	// Load asset from file system if not found in cache
	if _, ok := c.cache[assetType][name]; !ok {
		filepath := filepath.Join(c.path, assetPath[assetType], name)
		log.Printf("Loading asset %q", filepath)
		file, err := c.fs.Open(filepath)
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
