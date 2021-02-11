// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
	"time"
)

func main() {
	js.Global().Set("testCallback", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		callback := args[0]

		a := 0
		go func() {
			for {
				time.Sleep(5 * time.Second)
				a++
				callback.Invoke(fmt.Sprintf("Hello %d", a))
			}
		}()

		return js.Undefined()
	}))

	<-make(chan bool)
}
