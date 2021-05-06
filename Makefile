SHELL :=/bin/bash

all: server

bin/protoc:
	wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip && unzip -q -o protoc-3.15.8-linux-x86_64.zip

bin/protoc-gen-go:
	export GOBIN=$(shell pwd)/bin && \
		go install google.golang.org/protobuf/cmd/protoc-gen-go

bin/protoc-gen-go-grpc:
	export GOBIN=$(shell pwd)/bin && \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
bin/protoc-gen-grpc-gateway:
	export GOBIN=$(shell pwd)/bin && \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

bin/protoc-gen-openapiv2:
	export GOBIN=$(shell pwd)/bin && \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

bin/proctoc-gen-swagger:
	export GOBIN=$(shell pwd)/bin && \
		go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger


pre: bin/protoc-gen-grpc-gateway bin/protoc-gen-go-grpc bin/protoc-gen-openapiv2 bin/protoc-gen-go
	export GOBIN=$(shell pwd)/bin && \
	mkdir -p bin && \
	go version && \
	go mod tidy 
		

proto/helloworld/*.go: proto/helloworld/helloworld.proto bin/protoc pre
	export PATH=$(shell pwd)/bin:${PATH} && \
		export GOBIN=$(shell pwd)/bin && \
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
	rm -rf main server proto/helloworld/*.go cmd/helloworld ./bin
.PHONY: clean
