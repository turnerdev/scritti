all: cli web

cli:
	go build -o scritti main.go

web:
	GOARCH=wasm GOOS=js go build -ldflags="-s -w" -o www/lib.wasm main-wasm.go
	cp -rf "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" www/
	
test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html