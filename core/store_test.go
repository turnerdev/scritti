package core

import (
	"testing"
)

func TestFileStore(t *testing.T) {
	componentName := "example"

	fakeFileSystem := StubFS{
		map[string]StubFile{
			"test": {componentName},
		},
	}

	fileStore := FileStore{
		path: "",
		fs:   fakeFileSystem,
		cache: map[AssetType]map[string]string{
			ComponentType: {},
			StyleType:     {},
		},
	}

	t.Run("Test cache miss", func(t *testing.T) {

		asset, err := fileStore.Get(ComponentType, "test")
		component := asset.(*Component)

		if err != nil {
			t.Errorf("Cache error %q", err)
		}

		if component.name != componentName {
			t.Errorf("Expected component name %q got %q", componentName, component.name)
		}

	})

}
