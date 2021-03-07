all: cli web

cli:
	go build -o scritti main.go

web:
	GOARCH=wasm GOOS=js go build -ldflags="-s -w" -o www/lib.wasm wasm/main.go
	cp -rf "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" www/
	
tinygo:
	tinygo build -o www/wasm.wasm -target wasm wasm/main.go
	cp -rf $(shell dirname $(shell which tinygo))/../share/tinygo/targets/wasm_exec.js www/

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html