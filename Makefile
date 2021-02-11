all: cli web

cli:
	go build -o scritti main.go

web:
	GOARCH=wasm GOOS=js go build -o www/lib.wasm wasm/main.go