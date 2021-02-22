package core

import (
	"bytes"
	"errors"
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

	t.Run("Render error", func(t *testing.T) {
		delete(assets, AssetKey{StyleType, "child"})

		b := new(bytes.Buffer)
		err := RenderComponent(b, assets[AssetKey{ComponentType, "main"}].(Component), fn)
		if err.Error() != "Asset not found" {
			t.Error(err)
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
