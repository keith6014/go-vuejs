SHELL :=/bin/bash

all: main server


helloworld/helloworld_grpc.pb.go helloworld/helloworld.pb.go: helloworld/helloworld.proto
	bin/protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$^

main: cmd/main.go
	go build $^

server: helloworld/helloworld_grpc.pb.go helloworld/helloworld.pb.go | cmd/server.go
	go build cmd/server.go

clean:
	rm -rf main server helloworld/helloworld.pb.go
.PHONY: clean
