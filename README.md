# go-vuejs
golang with vuejs

https://grpc.io/blog/coreos/
https://dev.to/toransahu/part-1-building-a-basic-microservice-with-grpc-using-golang-304d

```bash
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```

```bash
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
	```
