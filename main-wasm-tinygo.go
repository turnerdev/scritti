package main

import (
	core "scritti/core"
	filesystem "scritti/filesystem"
)

var fs filesystem.FileSystem = filesystem.NewMemoryFileSystem()
var store core.AssetStore = core.NewFileStore(fs, "")

func mafin() {
	println("test")
}

// func wasmgetcallback(data map[string]interface{})

//export wasmget
// func wasmget(key core.AssetKey) map[string]interface{} {

// 	go func() {
// 		asset, err := store.Get(key)
// 		if err != nil {
// 			panic(err)
// 		}

// 		switch v := asset.(type) {
// 		case core.Component:
// 			buffer := new(bytes.Buffer)
// 			err = core.RenderComponent(buffer, v, store.Get)
// 			if err != nil {
// 				panic(err)
// 			}
// 			wasmgetcallback(map[string]interface{}{
// 				"id":     key,
// 				"html":   buffer.String(),
// 				"source": v.Source,
// 			})
// 		case core.Style:
// 			wasmgetcallback(map[string]interface{}{
// 				"id":     key,
// 				"source": v.Source,
// 			})
// 		}
// 	}()

// 	return nil
// }

//exports wasmset
func wasmset(params map[string]interface{}) {
	key := params["id"].(core.AssetKey)
	source := params["source"].(string)
	store.Set(key, source)
}
