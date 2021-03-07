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

func (s StubAssetStore) Get(key core.AssetKey) (core.Asset, error) {
	return s.component, nil
}

func (s StubAssetStore) Watch(key core.AssetKey, done <-chan bool) <-chan core.AssetEvent {
	ch := make(chan core.AssetEvent)
	return ch
}

func (s StubAssetStore) Close() error {
	return nil
}

func TestServer(t *testing.T) {
	store := StubAssetStore{
		&core.Component{},
	}

	server := NewComponentServer(&store)

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
