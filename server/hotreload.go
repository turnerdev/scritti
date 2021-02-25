package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	core "scritti/core"
	"strings"

	"golang.org/x/net/websocket"
)

var connectionUpgradeRegex = regexp.MustCompile("(^|.*,\\s*)upgrade($|\\s*,)")

func isWebsocketRequest(req *http.Request) bool {
	return connectionUpgradeRegex.MatchString(strings.ToLower(req.Header.Get("Connection"))) && strings.ToLower(req.Header.Get("Upgrade")) == "websocket"
}

type AssetData struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	HTML   string `json:"html"`
}

func (p ComponentServer) handleWebSockets(ws *websocket.Conn) {
	done := make(chan bool)
	key := core.AssetKey{AssetType: core.ComponentType, Name: "main"}

	log.Println("Connection established by " + ws.RemoteAddr().String())

	// Push loop
	go func() {
		for range p.store.Watch(key, done) {
			log.Println("hot reloading!")
			asset, err := p.store.Get(key)
			if err != nil {
				log.Fatal(err)
			}
			component := asset.(core.Component)

			// Render output
			buffer := new(bytes.Buffer)
			err = core.RenderComponent(buffer, asset.(core.Component), p.store.Get)
			if err != nil {
				log.Fatal(err)
			}

			data := AssetData{
				ID:     "main",
				Source: component.Source,
				HTML:   buffer.String(),
			}

			err = websocket.JSON.Send(ws, data)
			if err != nil {
				log.Println("message not sent " + err.Error())
				break
			}
			log.Println("done reload")
		}
	}()

	// RPC loop
	for {
		var request JsonRpcRequest
		err := websocket.JSON.Receive(ws, &request)

		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by " + ws.RemoteAddr().String())
				break
			}
			log.Println("Message not received " + err.Error())
			break
		}

		switch request.Method {
		case "get":
			var key core.AssetKey
			json.Unmarshal([]byte(request.Params), &key)
			log.Println(key)
			asset, err := p.store.Get(key)
			if err != nil {
				websocket.JSON.Send(ws, &JsonRpcResponse{
					JSONRPC: "2.0",
					Error: &JsonRpcError{
						Code:    1,
						Message: err.Error(),
					},
					ID: request.ID,
				})
			}

			buffer := new(bytes.Buffer)
			var response JsonRpcResponse

			switch v := asset.(type) {
			case core.Component:
				err = core.RenderComponent(buffer, v, p.store.Get)
				if err != nil {
					log.Fatal(err)
				}
				response = JsonRpcResponse{
					JSONRPC: "2.0",
					Result: &AssetData{
						ID:     key.Name,
						Source: v.Source,
						HTML:   buffer.String(),
					},
					ID: request.ID,
				}
			case core.Style:
				response = JsonRpcResponse{
					JSONRPC: "2.0",
					Result: &AssetData{
						ID:     key.Name,
						Source: v.Source,
					},
					ID: request.ID,
				}
			default:
				response = JsonRpcResponse{
					JSONRPC: "2.0",
					Error: &JsonRpcError{
						Code:    1,
						Message: fmt.Sprintf("Can't convert %d %q", key.AssetType, key.Name),
					},
					ID: request.ID,
				}
			}

			// Render output
			err = websocket.JSON.Send(ws, &response)
			if err != nil {
				log.Println("message not sent " + err.Error())
				break
			}
		default:
			log.Println("Unknown request method", request)
		}
	}

	log.Println("Done!!!")
	close(done)
}

// HotReload web socket handler
// func (p ComponentServer) handleWebSockets(ws *websocket.Conn) {
// 	m := Message{
// 		Message: "reload",
// 	}

// 	for range p.reload {
// 		err := websocket.JSON.Send(ws, &m)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		log.Println("Hot reloaded")
// 	}
// }

func (p ComponentServer) HandleHotReload(w http.ResponseWriter, r *http.Request) {
	if isWebsocketRequest(r) {
		websocket.Handler(p.handleWebSockets).ServeHTTP(w, r)
	}
}
