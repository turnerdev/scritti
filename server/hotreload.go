package server

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/websocket"
)

var connectionUpgradeRegex = regexp.MustCompile("(^|.*,\\s*)upgrade($|\\s*,)")

type Message struct {
	Message string
}

func isWebsocketRequest(req *http.Request) bool {
	return connectionUpgradeRegex.MatchString(strings.ToLower(req.Header.Get("Connection"))) && strings.ToLower(req.Header.Get("Upgrade")) == "websocket"
}

// HotReload web socket handler
func (p ComponentServer) handleWebSockets(ws *websocket.Conn) {
	m := Message{
		Message: "reload",
	}

	for range p.reload {
		err := websocket.JSON.Send(ws, &m)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Hot reloaded")
	}
}

func (p ComponentServer) HandleHotReload(w http.ResponseWriter, r *http.Request) {
	if isWebsocketRequest(r) {
		websocket.Handler(p.handleWebSockets).ServeHTTP(w, r)
	}
}
