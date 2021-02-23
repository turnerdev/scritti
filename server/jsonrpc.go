package server

import "encoding/json"

type JsonRpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      int             `json:"id"`
}

type JsonRpcResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result"`
	Error   *JsonRpcError `json:"error"`
	ID      int           `json:"id"`
}

type JsonRpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
