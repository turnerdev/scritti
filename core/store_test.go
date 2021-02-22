package core

import (
	"scritti/filesystem"
	"sync"
	"testing"
)

func TestFileStoreGet(t *testing.T) {
	fs := filesystem.NewMemoryFileSystem()
	fs.Write("main", "root\n\tnode1\n\tnode2")
	fs.Write("style/root", "class1\nclass2")
	fs.Write("style/node1", "class1\nclass2")
	fs.Write("style/node2", "class1\nclass2")
	store := NewFileStore(fs, "")

	t.Run("Test retrieve Component", func(t *testing.T) {
		asset, err := store.Get(AssetKey{ComponentType, "main"})
		if err != nil {
			t.Error(err)
		}

		component, ok := asset.(Component)
		if !ok {
			t.Errorf("Got %T, expected Component", asset)
		}
		if len(component.children) != 2 {
			t.Errorf("Got %d child elements, want 2", len(component.children))
		}
	})

	t.Run("Test retrieve Style", func(t *testing.T) {
		asset, err := store.Get(AssetKey{StyleType, "root"})
		if err != nil {
			t.Error(err)
		}

		style, ok := asset.(Style)
		if !ok {
			t.Errorf("Got %T, expected Style", asset)
		}
		if len(style.classes) != 2 {
			t.Errorf("Got %d child elements, want 2", len(style.classes))
		}
	})
}

func TestFileStoreWatch(t *testing.T) {
	fs := filesystem.NewMemoryFileSystem()
	fs.Write("main", "root\n\tnode1\n\tnode2")
	fs.Write("style/root", "class1\nclass2")
	fs.Write("style/node1", "class1\nclass2")
	fs.Write("style/node2", "class1\nclass2")

	t.Run("Test basic watch", func(t *testing.T) {
		store := NewFileStore(fs, "")
		done := make(chan bool)
		watch := store.Watch(AssetKey{ComponentType, "main"}, done)

		// Watch for 2 changes
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			<-watch
			<-watch
			wg.Done()
		}()

		// Make 2 changes
		fs.Write("main", "root\n\tnode1\n\tnode2\n\tnode3")
		fs.Write("main", "root\n\tnode1\n\tnode2\n\tnode4")

		wg.Wait()
		close(done)
	})

	t.Run("Test watch with dependencies", func(t *testing.T) {
		store := NewFileStore(fs, "")
		defer store.Close()
		done := make(chan bool)
		watch := store.Watch(AssetKey{ComponentType, "main"}, done)

		// Watch for 2 changes
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			<-watch
			wg.Done()
		}()

		// Make 2 changes
		fs.Write("style/node1", "class4")

		wg.Wait()
		close(done)
	})

}
