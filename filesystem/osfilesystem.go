// +build !js

package filesystem

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

// OSFileSystem implements FileSystem, exposing the OS file system
type OSFileSystem struct{}

func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

func (OSFileSystem) Open(name string) (File, error)        { return os.Open(name) }
func (OSFileSystem) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }
func (fs OSFileSystem) Watch(name string, done <-chan bool) (<-chan File, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	files := make(chan File)

	go func() {
		for {
			select {
			case <-done:
				watcher.Close()
				close(files)
				return
			case event := <-watcher.Events:
				log.Println("File system event:", name, event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("File system write:", name)
					file, err := fs.Open(name)
					if err != nil {
						log.Fatal(err)
					}
					files <- file
				}

			}
		}
	}()

	err = watcher.Add(name)
	if err != nil {
		log.Fatal(err)
	}

	return files, nil
}
