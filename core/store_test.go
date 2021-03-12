package core

import (
	"bufio"
	"scritti/filesystem"
	"sort"
	"sync"
	"testing"
)

func TestFileStoreGet(t *testing.T) {
	fs := filesystem.NewMemoryFileSystem()
	fsWrite(fs, "main", "root\n\tnode1\n\tnode2")
	fsWrite(fs, "style/root", "class1\nclass2")
	fsWrite(fs, "style/node1", "class1\nclass2")
	fsWrite(fs, "style/node2", "class1\nclass2")
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

func TestFileStoreSet(t *testing.T) {
	fs := filesystem.NewMemoryFileSystem()
	store := NewFileStore(fs, "")

	t.Run("Test update Component", func(t *testing.T) {
		key := AssetKey{ComponentType, "main"}

		store.Set(key, "a")
		store.Set(key, "ab")
		store.Set(key, "abc")
		store.Set(key, "abcd")

		asset, err := store.Get(AssetKey{ComponentType, "main"})
		if err != nil {
			t.Error(err)
		}

		component, ok := asset.(Component)
		if !ok {
			t.Errorf("Got %T, expected Component", asset)
		}

		if component.Source != "abcd" {
			t.Errorf("Got %q, want %q", component.Source, "abcd")
		}
	})
}

func fsWrite(fs filesystem.FileSystem, name string, content string) {
	file, err := fs.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.WriteString(content)
	w.Flush()
}

func TestFileStoreWatch(t *testing.T) {
	fs := filesystem.NewMemoryFileSystem()
	fsWrite(fs, "main", "root\n\tnode1\n\t\tnode2\n\t\tnode2")
	fsWrite(fs, "style/root", "class1\nclass2")
	fsWrite(fs, "style/node1", "class1\nclass2")
	fsWrite(fs, "style/node2", "class1\nclass2")

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
		fsWrite(fs, "main", "root\n\tnode1\n\tnode2\n\tnode3")
		fsWrite(fs, "main", "root\n\tnode1\n\tnode2\n\tnode4")

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
		fsWrite(fs, "style/node2", "class4")

		wg.Wait()
		close(done)
	})

}

func TestFileStoreList(t *testing.T) {
	fs := filesystem.NewMemoryFileSystem()
	fsWrite(fs, "main", "root\n\tnode1\n\t\tnode2\n\t\tnode2")
	fsWrite(fs, "style/root", "class1\nclass2")
	fsWrite(fs, "style/node1", "class1\nclass2")
	fsWrite(fs, "style/node2", "class1\nclass2")

	t.Run("Test list assets", func(t *testing.T) {
		store := NewFileStore(fs, "")
		defer store.Close()

		_, err := store.Get(AssetKey{ComponentType, "main"})
		if err != nil {
			t.Error(err)
		}

		want := [4]AssetKey{
			{ComponentType, "main"},
			{StyleType, "node1"},
			{StyleType, "node2"},
			{StyleType, "root"},
		}
		var got [4]AssetKey
		data := store.List()
		sort.Sort(ByAssetKey(data))
		copy(got[:], data)

		if want != got {
			t.Errorf("Got %q, want %q", got, want)
		}
	})
}

type ByAssetKey []AssetKey

func (a ByAssetKey) Len() int { return len(a) }
func (a ByAssetKey) Less(i, j int) bool {
	if a[i].AssetType > a[j].AssetType {
		return false
	} else if a[i].AssetType < a[j].AssetType {
		return true
	}
	return a[i].Name < a[j].Name
}
func (a ByAssetKey) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
