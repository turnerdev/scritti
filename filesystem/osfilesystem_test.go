package filesystem

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestOSFileSystem(t *testing.T) {

	t.Run("Test OS File System - Open", func(t *testing.T) {
		want := "content"
		filename := tempMkFile(t, "", "testfile")
		fs := NewOSFileSystem()

		WriteFile(t, filename, want)
		file, err := fs.Open(filename)
		if err != nil {
			t.Error(err)
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			t.Error(err)
		}

		got := strings.TrimSpace(string(data))

		if got != want {
			t.Errorf("Expected %q got %q", got, want)
		}

	})

	t.Run("Test OS File System - Watch", func(t *testing.T) {
		content := "old content"
		newContent := "new content"
		filename := tempMkFile(t, "", "testfile")
		fs := NewOSFileSystem()
		done := make(chan bool)

		// Write initial file value
		WriteFile(t, filename, content)

		// Watch file for modifications, receiving events
		events, _ := fs.Watch(filename, done)

		// Modify file
		WriteFile(t, filename, newContent)

		// Receive modification event
		_, ok := <-events
		if !ok {
			t.Error("Channel is closed")
		}

		// Close done channel to indicate we have finished watching
		close(done)

		// Expect watcher to close events channel
		_, ok = <-events
		if ok {
			t.Error("Channel is not closed")
		}

	})

}

func tempMkFile(t *testing.T, dir string, filename string) string {
	f, err := ioutil.TempFile(dir, filename)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer f.Close()
	return f.Name()
}

func WriteFile(t *testing.T, filename string, body string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 777)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(f, body)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
