package core

import (
	"bytes"
	"errors"
	"scritti/filesystem"
	"testing"
	// "golang.org/x/net/html"
)

func TestRender(t *testing.T) {
	assets := map[AssetKey]Asset{
		{ComponentType, "main"}: MakeComponent("parent\n\tchild"),
		{StyleType, "parent"}:   MakeStyle("one\ntwo"),
		{StyleType, "child"}:    MakeStyle("three\nfour"),
	}

	fn := func(assetKey AssetKey) (Asset, error) {
		if _, ok := assets[assetKey]; !ok {
			return struct{}{}, errors.New("Asset not found")
		}
		return assets[assetKey], nil
	}

	t.Run("Basic rendering", func(t *testing.T) {
		want := "<div class=\"one two\"><div class=\"three four\"></div></div>"

		b := new(bytes.Buffer)
		err := RenderComponent(b, assets[AssetKey{ComponentType, "main"}].(Component), fn)
		if err != nil {
			t.Error(err)
		}

		if got := b.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Render with missing styles", func(t *testing.T) {
		delete(assets, AssetKey{StyleType, "child"})
		want := "<div class=\"one two\"><div class=\"\"></div></div>"

		b := new(bytes.Buffer)
		err := RenderComponent(b, assets[AssetKey{ComponentType, "main"}].(Component), fn)
		if err != nil {
			t.Error(err)
		}

		if got := b.String(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func MakeComponent(source string) Component {
	component, err := NewComponent(source)
	if err != nil {
		panic(err)
	}
	return component
}

func MakeStyle(source string) Style {
	style, err := NewStyle(source)
	if err != nil {
		panic(err)
	}
	return style
}

func TestSampleData(t *testing.T) {
	fs := filesystem.NewOSFileSystem()
	store := NewFileStore(fs, "../sampledata")
	defer store.Close()

	asset, err := store.Get(AssetKey{ComponentType, "main"})
	if err != nil {
		t.Error(err)
	}
	b := new(bytes.Buffer)
	RenderComponent(b, asset.(Component), store.Get)
	t.Error("test")
}
