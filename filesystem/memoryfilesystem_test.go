package filesystem

import (
	"io/ioutil"
	"sync"
	"testing"
)

func TestMemoryFileSystem(t *testing.T) {
	filename := "testfile"

	t.Run("Test Memory File System - Open", func(t *testing.T) {
		want := "File content"
		fs := NewMemoryFileSystem()
		fs.Write(filename, want)
		file, _ := fs.Open(filename)
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
		fs.Write(filename, "Initial content")
		done := make(chan bool)

		// Initialize 2 watchers
		watch1, _ := fs.Watch(filename, done)
		watch2, _ := fs.Watch(filename, done)
		watchers := []<-chan File{watch1, watch2}

		var wg sync.WaitGroup
		wg.Add(len(watchers))

		// Goroutine to process 2 file system notifications
		processFn := func(i int) {
			defer wg.Done()
			data, _ := ioutil.ReadAll(<-watchers[i])
			got := string(data)
			if want[0] != got {
				t.Errorf("Want %q got %q", want[0], got)
			}
			data, _ = ioutil.ReadAll(<-watchers[i])
			got = string(data)
			if want[1] != got {
				t.Errorf("Want %q got %q", want[1], got)
			}
		}

		for i := range watchers {
			go processFn(i)
		}

		// Write 2 changes
		for _, data := range want {
			fs.Write(filename, data)
		}

		// Expect all goroutines to complete
		wg.Wait()

		// Close global channel
		close(done)

		// Ensure all watchers have been removed - no blocking channels
		fs.Write(filename, "test")
	})

	t.Run("Test Memory File System - Watch missing file", func(t *testing.T) {
		fs := NewMemoryFileSystem()
		fs.Write(filename, "Initial content")
		done := make(chan bool)

		_, err := fs.Watch("wrongfile", done)
		if err == nil {
			t.Fatal("Expected error")
		}

		close(done)
	})

}
