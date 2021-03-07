// +build js,wasm

package main

import (
	"bytes"
	"scritti/core"
	"scritti/filesystem"
	"syscall/js"
)

func main() {

	console := js.Global().Get("console")
	console.Call("log", "Hi!")

	fs := filesystem.NewMemoryFileSystem()
	store := core.NewFileStore(fs, "")
	defer store.Close()

	js.Global().Set("wasmget", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := args[0]
		if params.Type() != js.TypeObject {
			return js.Error{js.ValueOf("Invalid parameter")}
		}

		cost := core.AssetKey{
			core.AssetType(params.Get("assetType").Int()),
			params.Get("name").String(),
		}

		asset, err := store.Get(cost)
		if err != nil {
			return js.Error{js.ValueOf(err.Error())}
		}

		var result map[string]interface{}

		switch v := asset.(type) {
		case core.Component:
			buffer := new(bytes.Buffer)
			err = core.RenderComponent(buffer, v, store.Get)
			if err != nil {
				return js.Error{js.ValueOf(err.Error())}
			}
			result = map[string]interface{}{
				"id":     params,
				"html":   buffer.String(),
				"source": v.Source,
			}
		case core.Style:
			println("View source")
			println(v.Source)
			result = map[string]interface{}{
				"id":     params,
				"source": v.Source,
			}
		default:
			result = map[string]interface{}{}
		}

		return js.ValueOf(result)
	}))

	js.Global().Set("wasmset", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		params := args[0]
		if params.Type() != js.TypeObject {
			return js.Error{js.ValueOf("Invalid parameter")}
		}

		key := params.Get("id")

		cost := core.AssetKey{
			core.AssetType(key.Get("assetType").Int()),
			key.Get("name").String(),
		}

		println(cost.AssetType)
		println(cost.Name)

		err := store.Set(cost, params.Get("source").String())
		if err != nil {
			// console.Call("error", err.Error())
			// return js.Error{js.ValueOf(err.Error())}
		}

		return js.Undefined()
	}))

	<-make(chan bool)
}
