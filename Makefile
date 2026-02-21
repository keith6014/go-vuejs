SHELL :=/bin/bash

all: server

PROJECT_DIR := $(shell pwd)
export GOMODCACHE := $(PROJECT_DIR)/.gocache
export GOBIN := $(PROJECT_DIR)/bin
export PATH := $(PROJECT_DIR)/bin:$(PATH)

bin/protoc:
	wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip && unzip -q -o protoc-3.15.8-linux-x86_64.zip

pre:
	mkdir -p $(PROJECT_DIR)/bin && \
	go mod tidy && \
	go install  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	go install github.com/jstemmer/go-junit-report/v2@latest


proto/helloworld/*.go: proto/helloworld/helloworld.proto bin/protoc pre
	bin/protoc -I ./proto \
		-I $(GOMODCACHE)/github.com/grpc-ecosystem/grpc-gateway/v2@v2.4.0/ \
		--go_out ./proto --go_opt paths=source_relative \
		--go-grpc_out ./proto --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
		--swagger_out=logtostderr=true:cmd --go-grpc_opt paths=source_relative \
		proto/helloworld/helloworld.proto


test: cmd/server.go
	     go test ./... -v | tee >(go-junit-report > report.xml)

server: proto/helloworld/*.go cmd/server.go 
	     go build cmd/server.go

clean:
	rm -rf main server proto/helloworld/*.go cmd/helloworld
.PHONY: clean
