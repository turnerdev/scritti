package core

import (
	"errors"
	"io/ioutil"
	"log"
	"path/filepath"
	"scritti/filesystem"
	"sync"
)

// AssetType enum
type AssetType int

// AssetType enum
const (
	ComponentType AssetType = iota
	StyleType
)

// AssetEvent represents change event emitted from by an Asset Store
type AssetEvent struct {
	key    AssetKey
	direct bool
}

// AssetKey is a composite key for Assets in the store
type AssetKey struct {
	assetType AssetType
	name      string
}

// assetEntry is the internal representation of an Asset in the store
type assetEntry struct {
	asset        Asset
	dependencies []AssetKey
	dependants   map[AssetKey]struct{}
	mu           sync.RWMutex
	watchers     map[chan AssetEvent]struct{}
}

// newAssetEntry returns a pointer to a new AssetValue instance
func newAssetEntry(asset Asset) *assetEntry {
	return &assetEntry{
		asset:        asset,
		dependencies: []AssetKey{},
		dependants:   make(map[AssetKey]struct{}),
		mu:           sync.RWMutex{},
		watchers:     make(map[chan AssetEvent]struct{}),
	}
}

// AssetStore provides an interface to retrieve components
type AssetStore interface {
	Get(key AssetKey) (Asset, error)
	Watch(key AssetKey, done <-chan bool) <-chan AssetEvent
	Close() error
}

// FileStore TODO
type FileStore struct {
	path    string
	fs      filesystem.FileSystem
	entries map[AssetKey]*assetEntry
	done    chan bool
}

// NewFileStore returns a new File Store
func NewFileStore(fs filesystem.FileSystem, path string) *FileStore {
	return &FileStore{
		path:    path,
		fs:      fs,
		entries: make(map[AssetKey]*assetEntry),
		done:    make(chan bool),
	}
}

var assetPath = map[AssetType]string{
	StyleType: "style",
}

// updateAsset refreshes an entry in the store from the filesystem
func (c FileStore) updateAsset(key AssetKey) error {
	// Get current asset state
	// current, ok := c.entries[key]
	// if !ok {
	// 	return fmt.Errorf("Cannot find existing Asset entry for %q", key.name)
	// }
	// // Load new asset from
	// next, err := c.loadAsset(key)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// fetchAsset retrieves an Asset from the file system
func (c FileStore) fetchAsset(key AssetKey) (Asset, error) {
	// Open asset source in file system
	path := c.getPath(key)
	log.Printf("Loading asset %q", path)
	file, err := c.fs.Open(path)
	if err != nil {
		log.Printf("%+v", errors.New("Asset not found"))
		return nil, err
	}

	// Read asset source to buffer
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	source := string(data)

	// Ensure Asset compiles
	asset, err := NewAssetFactory(key.assetType, source)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// getAsset returns the internal store representation of an Asset. I
func (c FileStore) getAsset(key AssetKey) (*assetEntry, error) {
	if _, ok := c.entries[key]; !ok {
		err := c.addEntry(key)
		if err != nil {
			return nil, err
		}
	}
	return c.entries[key], nil
}

// Compare old/new
// Update dependants
func (c FileStore) updateAssetDependencies(key AssetKey) error {
	diff := map[AssetKey]int{}
	asset, err := c.getAsset(key)
	if err != nil {
		return err
	}

	// Update dependency map
	// -1: Dependency exists only in previous state
	//  0: Dependency exists in both previous and new state
	//  1: Dependency exists only in new state
	for _, oldKey := range asset.dependencies {
		diff[oldKey] = -1
	}

	// Get new dependencies
	newDependencies := getDependencyKeys(asset.asset)

	for _, newKey := range newDependencies {
		if _, ok := diff[newKey]; !ok {
			diff[newKey] = 1
		} else {
			diff[newKey] = 0
		}
	}

	// Update asset to new dependencies
	asset.mu.Lock()
	asset.dependencies = newDependencies
	asset.mu.Unlock()

	// Update dependants
	for k, v := range diff {
		dependency, err := c.getAsset(k)
		if err != nil {
			return err
		} else if v == -1 {
			// Remove old dependant
			delete(dependency.dependants, key)
		} else if v == 1 {
			// Remove old dependant
			dependency.dependants[key] = struct{}{}
		}
	}

	return nil
}

func (c FileStore) addEntry(key AssetKey) error {
	// Load Asset from file system
	asset, err := c.fetchAsset(key)
	if err != nil {
		return err
	}
	// Create new store entry for Asset
	assetEntry := newAssetEntry(asset)
	c.entries[key] = assetEntry
	err = c.updateAssetDependencies(key)
	if err != nil {
		return err
	}

	// Create watcher channel
	assetWatch, err := c.fs.Watch(c.getPath(key), c.done)
	if err != nil {
		return err
	}

	go func() {
		for range assetWatch {
			// Notify any watchers subscribed to the asset
			assetEntry.mu.RLock()
			assetEntry.asset, err = c.fetchAsset(key)
			if err != nil {
				panic(err)
			}
			for watcher := range assetEntry.watchers {
				watcher <- AssetEvent{key, true}
			}
			assetEntry.mu.RUnlock()
		}
	}()

	return nil
}

func (c FileStore) getPath(key AssetKey) string {
	return filepath.Join(c.path, assetPath[key.assetType], key.name)
}

// Watch an Asset in the store, subscribing to changes
func (c FileStore) Watch(key AssetKey, done <-chan bool) <-chan AssetEvent {
	ch := make(chan AssetEvent)

	_, err := c.Get(key)
	if err != nil {
		log.Fatalf("Unable to watch %q, %v", key.name, err)
	}

	// Subscribe to change events for asset
	assetEntry := c.entries[key]
	assetEntry.mu.Lock()
	assetEntry.watchers[ch] = struct{}{}
	assetEntry.mu.Unlock()

	// depWatchers := []<-chan AssetEvent{}
	// for _, dependency := range getDependencyKeys(asset) {
	// 	dependency = c.entries
	// 	depWatchers = append(depWatchers, w)
	// }
	// depWatcher := merge(ch, depWatchers...)

	go func() {
		<-done
		// Clear channel before closing
		go func() {
			for range ch {
			}
		}()
		assetEntry.mu.Lock()
		delete(assetEntry.watchers, ch)
		assetEntry.mu.Unlock()
		close(ch)
	}()

	return ch
}

// Get TODO
func (c FileStore) Get(key AssetKey) (Asset, error) {
	// Load asset from file system if not found in cache
	asset, err := c.getAsset(key)
	if err != nil {
		return nil, err
	}
	return asset.asset, nil
}

// Close all channels and file system watches
func (c FileStore) Close() error {
	close(c.done)
	return nil
}

func merge(ins ...<-chan filesystem.File) <-chan filesystem.File {
	out := make(chan filesystem.File)
	var wg sync.WaitGroup
	wg.Add(len(ins))
	for _, in := range ins {
		go func(in <-chan filesystem.File) {
			for event := range in {
				out <- event
			}
			wg.Done()
		}(in)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func indexOf(slice []AssetKey, element AssetKey) int {
	for i, item := range slice {
		if item == element {
			return i
		}
	}
	return -1
}

func removeAt(s []AssetKey, i int) []AssetKey {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
