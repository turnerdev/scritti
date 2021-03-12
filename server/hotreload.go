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
	ID     core.AssetKey `json:"id"`
	Source string        `json:"source"`
	HTML   string        `json:"html"`
}

// makeError returns the appropriate JSON RPC Error for an error type
func makeError(id int, err error) JsonRpcResponse {
	var errorDetail *JsonRpcError

	switch err.(type) {
	case *core.AssetNotFound:
		errorDetail = &JsonRpcError{
			Code:    1,
			Message: err.Error(),
		}
	default:
		errorDetail = &JsonRpcError{
			Code:    0,
			Message: err.Error(),
		}
	}

	return JsonRpcResponse{
		JSONRPC: "2.0",
		Error:   errorDetail,
		ID:      id,
	}
}

func (p ComponentServer) pushLoop(ws *websocket.Conn, done <-chan bool) {
	key := core.AssetKey{AssetType: core.ComponentType, Name: "main"}

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
			ID:     core.AssetKey{0, "main"},
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
}

func (p ComponentServer) setAction(request JsonRpcRequest) JsonRpcResponse {
	var data AssetData
	json.Unmarshal([]byte(request.Params), &data)
	log.Printf("Set: %q\n", data)
	err := p.store.Set(data.ID, data.Source)

	if err != nil {
		return JsonRpcResponse{
			JSONRPC: "2.0",
			Error: &JsonRpcError{
				Code:    1,
				Message: err.Error(),
			},
			ID: request.ID,
		}
	}

	return JsonRpcResponse{
		JSONRPC: "2.0",
		Result: &AssetData{
			ID:     data.ID,
			Source: data.Source,
			HTML:   data.HTML,
		},
		ID: request.ID,
	}
}

func (p ComponentServer) getAction(request JsonRpcRequest) JsonRpcResponse {
	var key core.AssetKey
	json.Unmarshal([]byte(request.Params), &key)
	log.Printf("Get: %q\n", key)

	buffer := new(bytes.Buffer)

	asset, err := p.store.Get(key)
	if err != nil {
		return makeError(request.ID, err)
	}

	switch v := asset.(type) {
	case core.Component:
		err = core.RenderComponent(buffer, v, p.store.Get)
		if err != nil {
			log.Fatal(err)
		}
		return JsonRpcResponse{
			JSONRPC: "2.0",
			Result: &AssetData{
				ID:     key,
				Source: v.Source,
				HTML:   buffer.String(),
			},
			ID: request.ID,
		}
	case core.Style:
		return JsonRpcResponse{
			JSONRPC: "2.0",
			Result: &AssetData{
				ID:     key,
				Source: v.Source,
			},
			ID: request.ID,
		}
	case core.SVG:
		return JsonRpcResponse{
			JSONRPC: "2.0",
			Result: &AssetData{
				ID:     key,
				Source: v.Source,
			},
			ID: request.ID,
		}
	}

	return JsonRpcResponse{
		JSONRPC: "2.0",
		Error: &JsonRpcError{
			Code:    1,
			Message: fmt.Sprintf("Can't convert %d %q", key.AssetType, key.Name),
		},
		ID: request.ID,
	}
}

func (p ComponentServer) listAction(request JsonRpcRequest) JsonRpcResponse {
	return JsonRpcResponse{
		JSONRPC: "2.0",
		Error: &JsonRpcError{
			Code:    1,
			Message: fmt.Sprintf("Can't convert %d %q", key.AssetType, key.Name),
		},
		ID: request.ID,
	}
}

func (p ComponentServer) rpcLoop(ws *websocket.Conn) {
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
		case "set":
			response := p.setAction(request)
			// Render output
			err = websocket.JSON.Send(ws, &response)
			if err != nil {
				log.Println("message not sent " + err.Error())
				break
			}

		case "get":

			var response *JsonRpcResponse

			// Render output
			err = websocket.JSON.Send(ws, &response)
			if err != nil {
				log.Println("Message not sent " + err.Error())
				break
			}
		default:
			log.Println("Unknown request method", request)
		}
	}

}

// handleWebSockets manages rpc calls and push notifications for a socket
func (p ComponentServer) handleWebSockets(ws *websocket.Conn) {
	done := make(chan bool)

	log.Println("Connection established by " + ws.RemoteAddr().String())

	// Push loop
	go p.pushLoop(ws, done)

	// RPC loop
	p.rpcLoop(ws)

	log.Println("Socket disconnected by user")
	close(done)
}

func (p ComponentServer) HandleHotReload(w http.ResponseWriter, r *http.Request) {
	if isWebsocketRequest(r) {
		websocket.Handler(p.handleWebSockets).ServeHTTP(w, r)
	}
}
