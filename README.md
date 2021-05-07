# go-vuejs
golang with vuejs

https://grpc.io/blog/coreos/
https://github.com/grpc/grpc-go/tree/master/examples/helloworld
https://dev.to/toransahu/part-1-building-a-basic-microservice-with-grpc-using-golang-304d

```bash
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```


### google apis ###
mkdir -p proto/google/api

```bash
wget https://github.com/googleapis/googleapis/archive/refs/heads/master.zip
unzip master.zip
mv googleapis-master/google/api/annotations.proto proto/google/api/
```
https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/adding_annotations/#using-protoc

### test ###
```bash
curl -X POST -k localhost:8080/hello_world -d '{"name":"yo"}'
```
