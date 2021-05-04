SHELL :=/bin/bash

all: server


proto/helloworld/*.go: proto/helloworld/helloworld.proto
	bin/protoc -I ./proto \
		--go_out ./proto --go_opt paths=source_relative \
		--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
		proto/helloworld/helloworld.proto

main: cmd/main.go
	go build $^

server: proto/helloworld/*.go | cmd/server.go
	go build cmd/server.go

clean:
	rm -rf main server proto/helloworld/*.go
.PHONY: clean
