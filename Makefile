SHELL :=/bin/bash

all: server


proto/helloworld/*.go: proto/helloworld/helloworld.proto
	bin/protoc -I ./proto \
		-I /home/user/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.4.0/ \
		--go_out ./proto --go_opt paths=source_relative \
		--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
		--swagger_out=logtostderr=true:cmd \
		proto/helloworld/helloworld.proto

main: cmd/main.go
	go build $^

server: proto/helloworld/*.go cmd/server.go
	go build cmd/server.go

clean:
	rm -rf main server proto/helloworld/*.go cmd/helloworld
.PHONY: clean
