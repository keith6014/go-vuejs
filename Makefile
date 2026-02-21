SHELL :=/bin/bash

all: server

PROJECT_DIR := $(shell pwd)
export GOMODCACHE := $(PROJECT_DIR)/.gocache
export GOBIN := $(PROJECT_DIR)/bin
export PATH := $(PROJECT_DIR)/bin:$(PATH)

bin/protoc:
	wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v33.5/protoc-33.5-linux-x86_64.zip && unzip -q -o protoc-33.5-linux-x86_64.zip && \
		touch $@

pre: 
	mkdir -p $(PROJECT_DIR)/bin && \
		go mod tidy && \
		go install  github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 && \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
		go install google.golang.org/protobuf/cmd/protoc-gen-go && \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
		go install github.com/jstemmer/go-junit-report/v2@latest


proto_generated/helloworld/*.go: bin/protoc | pre
	@echo "Testing Protoc" && \
		protoc --version && \
		mkdir -vp ./proto_generated && \
		protoc -I ./proto -I ./third_party \
		--go_out ./proto_generated --go_opt paths=source_relative \
		--go-grpc_out ./proto_generated --go-grpc_opt paths=source_relative \
		--grpc-gateway_out ./proto_generated --grpc-gateway_opt paths=source_relative \
		--openapiv2_out cmd/ \
		proto/helloworld/helloworld.proto  

test: server
	go test ./... -v | tee >(go-junit-report > report.xml)

server: proto_generated/helloworld/*.go cmd/server.go 
	go build cmd/server.go

clean:
	rm -rf main server proto/helloworld/*.go cmd/helloworld bin/* cmd/helloworld protoc-*.zip* proto_generated report.xml
.PHONY: clean
