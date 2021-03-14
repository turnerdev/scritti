#!/bin/sh
# docker run --rm -v $(go env GOPATH):/go -v $(go env GOROOT):/goroot -v $(pwd):/scritti  -e "GOPATH=/go" -e "GOROOT=/goroot" -e "GO111MODULE=on" tinygo/tinygo:0.17.0 tinygo build -o /scritti/wasm.wasm -target wasm --no-debug /scritti/main-wasm.go
# docker run --name tinygo_test --rm -i -t -v $(go env GOPATH):/go -v $(go env GOROOT):/goroot -v $(pwd):/scritti  -e "GOPATH=/go" -e "GOROOT=/goroot" -e "GO111MODULE=on" tinygo/tinygo:0.17.0 bash
tinygo build -o /scritti/www/tiny-wasm.wasm -target wasm --no-debug /scritti/main-wasm.go

