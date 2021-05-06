SHELL :=/bin/bash

all: server

bin/protoc:
	wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip && unzip -q -o protoc-3.15.8-linux-x86_64.zip

pre:
	export GOBIN=$(shell pwd)/bin && \
	mkdir -p bin && \
	go mod tidy && \
	go install \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
	        google.golang.org/protobuf/cmd/protoc-gen-go \
		    google.golang.org/grpc/cmd/protoc-gen-go-grpc
		

proto/helloworld/*.go: proto/helloworld/helloworld.proto bin/protoc pre
	export GOBIN=$(shell pwd)/bin && export PATH=${GOBIN}:${PATH} && \
	bin/protoc -I ./proto \
		-I $(shell go env GOMODCACHE)/github.com/grpc-ecosystem/grpc-gateway/v2@v2.4.0/ \
		--go_out ./proto --go_opt paths=source_relative \
		--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
		--swagger_out=logtostderr=true:cmd --go-grpc_opt paths=source_relative \
		proto/helloworld/helloworld.proto

main: cmd/main.go
	go build $^

server: proto/helloworld/*.go cmd/server.go 
	go build cmd/server.go

clean:
	rm -rf main server proto/helloworld/*.go cmd/helloworld
.PHONY: clean
