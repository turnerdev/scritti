package server

import (
	"net/http"
	"net/http/httptest"
	"scritti/core"
	"testing"
)

type StubAssetStore struct {
	component *core.Component
}

func (s StubAssetStore) Get(assetType core.AssetType, name string) (core.Asset, error) {
	return s.component, nil
}

func (s StubAssetStore) Watch(assetType core.AssetType, name string, ch chan<- bool) {

}

func (s StubAssetStore) Clear(assetType core.AssetType, name string) {

}

func TestServer(t *testing.T) {
	store := StubAssetStore{
		&core.Component{},
	}
	server := &ComponentServer{
		reload: make(chan bool),
		store:  &store,
	}

	t.Run("Test simple component rendering", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()

		if response.Result().StatusCode != 200 {
			t.Errorf("got %q", got)
		}
	})
}
