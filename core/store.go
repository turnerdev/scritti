package core

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"scritti/filesystem"
	"sync"
)

type AssetNotFound struct {
	asset AssetKey
}

func (e *AssetNotFound) Error() string {
	return fmt.Sprintf("Asset not found [%d] %s", e.asset.AssetType, e.asset.Name)
}

// AssetType enum
type AssetType int

// AssetType enum
const (
	ComponentType AssetType = iota
	StyleType
	SVGType
)

// AssetStatus enum
type AssetStatus int

// AssetStatus enum
const (
	NotLoaded AssetStatus = iota
	Loaded
)

// AssetEvent represents change event emitted from by an Asset Store
type AssetEvent struct {
	key    AssetKey
	direct bool
}

// AssetKey is a composite key for Assets in the store
type AssetKey struct {
	AssetType AssetType `json:"assetType"`
	Name      string    `json:"name"`
}

// assetEntry is the internal representation of an Asset in the store
type assetEntry struct {
	asset        Asset
	dependencies []AssetKey
	dependants   map[AssetKey]struct{}
	mu           sync.RWMutex
	watchers     map[chan AssetEvent]struct{}
	status       AssetStatus
}

// newAssetEntry returns a pointer to a new AssetValue instance
func newAssetEntry() *assetEntry {
	return &assetEntry{
		dependencies: []AssetKey{},
		dependants:   make(map[AssetKey]struct{}),
		mu:           sync.RWMutex{},
		watchers:     make(map[chan AssetEvent]struct{}),
		status:       NotLoaded,
	}
}

// AssetStore provides an interface to retrieve components
type AssetStore interface {
	Set(key AssetKey, content string) error
	Get(key AssetKey) (Asset, error)
	List() []AssetKey
	Watch(key AssetKey, done <-chan bool) <-chan AssetEvent
	Close() error
}

// FileStore TODO
type FileStore struct {
	path    string
	fs      filesystem.FileSystem
	entries map[AssetKey]*assetEntry
	mu      sync.RWMutex
	done    chan bool
}

// NewFileStore returns a new File Store
func NewFileStore(fs filesystem.FileSystem, path string) *FileStore {
	return &FileStore{
		path:    path,
		fs:      fs,
		entries: make(map[AssetKey]*assetEntry),
		mu:      sync.RWMutex{},
		done:    make(chan bool),
	}
}

var assetPath = map[AssetType]string{
	StyleType: "style",
	SVGType:   "svg",
}

// fetchAsset retrieves an Asset from the file system
func (c *FileStore) fetchAsset(key AssetKey) (Asset, error) {
	// Open asset source in file system
	path := c.getPath(key)
	log.Printf("Loading asset %q", path)

	file, err := c.fs.Open(path)

	if err != nil {
		log.Println(err)
		switch err.(type) {
		case *filesystem.FileNotFound:
			return nil, &AssetNotFound{key}
		default:
			return nil, err
		}
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Read asset source to buffer
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	source := string(data)

	// Ensure Asset compiles
	asset, err := NewAssetFactory(key.AssetType, source)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

// getAssetEntry returns the internal store representation of an Asset.
// If no entry exists, a new entry is first created
func (c *FileStore) getAssetEntry(key AssetKey) (*assetEntry, error) {
	// Create Asset entry if non existant
	c.mu.RLock()
	assetEntry, ok := c.entries[key]
	c.mu.RUnlock()

	if !ok {
		c.mu.Lock()
		c.entries[key] = newAssetEntry()
		assetEntry = c.entries[key]
		c.mu.Unlock()
	}

	if assetEntry.status == NotLoaded {
		err := c.loadAssetEntry(key)
		// TODO: Check for AssetNotFound error
		if _, ok := err.(*AssetNotFound); err != nil && !ok {
			return nil, err
		}
	}

	return assetEntry, nil
}

func (c *FileStore) updateAssetEntry(key AssetKey) error {
	// Look up Asset entry
	assetEntry, err := c.getAssetEntry(key)
	if err != nil {
		return err
	}

	// Caclulate dependency changes between current and new asset state
	// -1: Dependency exists only in previous state
	//  0: Dependency exists in both previous and new state
	//  1: Dependency exists only in new state
	diff := map[AssetKey]int{}

	// Get current dependencies
	for _, oldKey := range assetEntry.dependencies {
		diff[oldKey] = -1
	}

	// Fetch latest Asset version from file system
	newAsset, err := c.fetchAsset(key)
	if err != nil {
		return err
	}

	// Get new dependencies
	newDependencies := getDependencyKeys(newAsset)
	for _, newKey := range newDependencies {
		if _, ok := diff[newKey]; ok {
			diff[newKey] = 0
		} else {
			diff[newKey] = 1
		}
	}

	// Update asset and dependencies
	assetEntry.mu.Lock()
	assetEntry.asset = newAsset
	assetEntry.dependencies = newDependencies
	assetEntry.mu.Unlock()

	var wg sync.WaitGroup
	wg.Add(len(diff))

	// Update dependants
	for k, v := range diff {
		go func(k AssetKey, v int) {
			dependency, err := c.getAssetEntry(k)
			if err != nil {
				log.Printf("Can't load dependency %q of %q (%s)\n", k.Name, key.Name, err)
			} else if v == -1 {
				// Remove old dependant
				delete(dependency.dependants, key)
			} else if v == 1 {
				// Add new dependant
				log.Printf("> Add new dependant %q of %q", k.Name, key.Name)
				log.Println(dependency, err)
				dependency.dependants[key] = struct{}{}
			}
			wg.Done()
		}(k, v)
	}

	wg.Wait()

	return nil
}

func (c *FileStore) notifyWatchers(key AssetKey) error {
	assetEntry, err := c.getAssetEntry(key)
	if err != nil {
		return nil
	}
	assetEntry.mu.RLock()
	for watcher := range assetEntry.watchers {
		log.Printf("Notifying watcher (%q)\n", key.Name)
		watcher <- AssetEvent{key, true}
		log.Println("**Notified")
	}

	log.Printf("Notifying %d dependants of %s\n", len(assetEntry.dependants), key.Name)
	for dependant := range assetEntry.dependants {
		log.Printf("Bubbling change event (%q)\n", dependant.Name)
		c.notifyWatchers(dependant)
	}
	assetEntry.mu.RUnlock()
	return nil
}

func (c *FileStore) loadAssetEntry(key AssetKey) error {
	log.Printf("Loading asset entry %d-%s", key.AssetType, key.Name)

	// Fetch latest Asset version from file system
	_, err := c.fetchAsset(key)
	if err != nil {
		return err
	}

	// Update Asset Entry status
	c.mu.RLock()
	c.entries[key].status = Loaded
	c.mu.RUnlock()

	// Update Asset from the file system
	err = c.updateAssetEntry(key)
	if err != nil {
		return err
	}

	// Watch for changes to source in the file system
	watch, err := c.fs.Watch(c.getPath(key), c.done)
	if err != nil {
		return err
	}

	// When source changes, update asset entry and notify subscribers
	go func() {
		for range watch {
			log.Printf("Detected change in %q\n", c.getPath(key))

			// Notify any watchers subscribed to the asset
			c.updateAssetEntry(key)
			c.notifyWatchers(key)
		}
	}()

	return nil
}

func (c *FileStore) getPath(key AssetKey) string {
	return filepath.Join(c.path, assetPath[key.AssetType], key.Name)
}

// Watch an Asset in the store, subscribing to changes
func (c *FileStore) Watch(key AssetKey, done <-chan bool) <-chan AssetEvent {
	ch := make(chan AssetEvent)

	_, err := c.Get(key)
	if err != nil {
		log.Fatalf("Unable to watch %q, %v", key.Name, err)
	}

	// Subscribe to change events for asset
	c.mu.RLock()
	assetEntry := c.entries[key]
	c.mu.RUnlock()
	assetEntry.mu.Lock()

	assetEntry.watchers[ch] = struct{}{}
	assetEntry.mu.Unlock()

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

// Set creates or updates an Asset in the store with the given content
func (c *FileStore) Set(key AssetKey, content string) error {
	path := c.getPath(key)

	file, err := c.fs.Create(path)
	if err != nil {
		log.Printf("%+v", errors.New("Asset not found"))
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	w := bufio.NewWriter(file)
	w.WriteString(content)
	w.Flush()

	_, err = c.getAssetEntry(key)
	if err != nil {
		return err
	}

	return nil
}

// Get returns and Asset from the store
func (c *FileStore) Get(key AssetKey) (Asset, error) {
	// Load asset from file system if not found in cache
	asset, err := c.getAssetEntry(key)
	if err != nil {
		return nil, err
	}

	if asset.status != Loaded {
		return nil, errors.New("Asset not loaded")
	}

	return asset.asset, nil
}

// List returns the AssetKeys of all entries in the store
func (c *FileStore) List() []AssetKey {
	result := make([]AssetKey, 0)
	for k := range c.entries {
		result = append(result, k)
	}
	return result
}

// Close all channels and file system watches
func (c *FileStore) Close() error {
	close(c.done)
	return nil
}
