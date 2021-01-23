package server

import (
	"net/http"
	"net/http/httptest"
	"scritti/core"
	"testing"
)

type StubComponentStore struct {
	component *core.Component
}

func (s *StubComponentStore) Get(name string) (*core.Component, error) {
	return s.component, nil
}

func TestServer(t *testing.T) {
	store := StubComponentStore{
		&core.Component{},
	}
	server := &ComponentServer{&store}

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
