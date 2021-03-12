package filesystem

import (
	"bufio"
	"io/ioutil"
	"sync"
	"testing"
)

func TestMemoryFileSystem(t *testing.T) {
	filename := "testfile"

	t.Run("Test Memory File System - Open", func(t *testing.T) {
		want := "File content"
		fs := NewMemoryFileSystem()
		fsWrite(fs, filename, want)
		file, _ := fs.Open(filename)
		data, _ := ioutil.ReadAll(file)
		got := string(data)

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("Test Memory File System - Sequential writes", func(t *testing.T) {
		want := "New content"
		fs := NewMemoryFileSystem()
		fsWrite(fs, filename, "Old content")
		fsWrite(fs, filename, want)
		file, _ := fs.Open(filename)
		data, _ := ioutil.ReadAll(file)
		got := string(data)

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("Test Memory File System - Open", func(t *testing.T) {
		want := "File content"
		fs := NewMemoryFileSystem()
		fsWrite(fs, filename, want)
		file, _ := fs.Open(filename)
		data, _ := ioutil.ReadAll(file)
		got := string(data)

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("Test Memory File System - Create", func(t *testing.T) {
		want := "File\nContent"
		fs := NewMemoryFileSystem()
		file, err := fs.Create(filename)
		if err != nil {
			t.Error(err)
		}

		fsWrite(fs, filename, want)

		file, err = fs.Open(filename)
		if err != nil {
			t.Error(err)
		}
		data, _ := ioutil.ReadAll(file)
		got := string(data)

		if got != want {
			t.Errorf("Got %q, want %q", got, want)
		}
	})

	t.Run("Test Memory File System - Open missing file", func(t *testing.T) {
		fs := NewMemoryFileSystem()
		_, err := fs.Open(filename)

		if err == nil {
			t.Errorf("Expected error")
		}
	})

	t.Run("Test Memory File System - Watch", func(t *testing.T) {
		want := []string{
			"First change",
			"Second change",
		}

		// Create memory file system with a single file
		fs := NewMemoryFileSystem()
		fsWrite(fs, filename, "Initial content")
		done := make(chan bool)

		// Initialize 2 watchers
		watch1, _ := fs.Watch(filename, done)
		watch2, _ := fs.Watch(filename, done)
		watchers := []<-chan bool{watch1, watch2}

		var wg sync.WaitGroup
		wg.Add(len(watchers))

		// Goroutine to process 2 file system notifications
		processFn := func(i int) {
			defer wg.Done()
			<-watchers[i]
			<-watchers[i]
			file, _ := fs.Open(filename)
			data, _ := ioutil.ReadAll(file)
			got := string(data)
			if want[1] != got {
				t.Errorf("Want %q got %q", want, got)
			}
		}

		for i := range watchers {
			go processFn(i)
		}

		// Write 2 changes
		for _, data := range want {
			fsWrite(fs, filename, data)
		}

		// Expect all goroutines to complete
		wg.Wait()

		// Close global channel
		close(done)

		// Ensure all watchers have been removed - no blocking channels
		fsWrite(fs, filename, "test")
	})

	t.Run("Test Memory File System - Watch missing file", func(t *testing.T) {
		fs := NewMemoryFileSystem()
		fsWrite(fs, filename, "Initial content")
		done := make(chan bool)

		_, err := fs.Watch("wrongfile", done)
		if err == nil {
			t.Fatal("Expected error")
		}

		close(done)
	})

}

func fsWrite(fs FileSystem, name string, content string) {
	file, err := fs.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.WriteString(content)
	w.Flush()
}
