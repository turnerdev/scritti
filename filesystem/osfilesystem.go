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

func (OSFileSystem) Create(name string) (File, error) { return os.Create(name) }

func (OSFileSystem) Open(name string) (File, error) {
	file, err := os.Open(name)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			return nil, &FileNotFound{name}
		default:
			return nil, err
		}
	}
	return file, nil
}
func (OSFileSystem) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

func (fs OSFileSystem) Watch(name string, done <-chan bool) (<-chan bool, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	files := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				watcher.Close()
				close(files)
				return
			case event := <-watcher.Events:
				log.Printf("event:%s,%s", event.Name, event.Op)
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					if FileExist(name) {
						log.Printf("Reread file:%s ", name)
						err := watcher.Add(name)
						if err != nil {
							log.Println("error:", err)
						}
					}
				}
				// only work on WRITE events of the original filename
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}

				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("File system write:", name)
					files <- true
				}

			}
		}
	}()

	err = watcher.Add(name)
	if err != nil {
		return nil, err
		// log.Fatal(err)
	}

	return files, nil
}

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
