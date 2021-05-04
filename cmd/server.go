package main

import (
	"context"
	"log"
	"net"
	"net/http"

	pb "go-vuejs/proto/helloworld"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: in.Name + " world"}, nil
}

func main() {
	log.Println("grpc server on", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		log.Fatal(s.Serve(lis))
	}()
	conn, err := grpc.DialContext(
		context.Background(),
		port,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("err ", err)
	}
	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8080",
		Handler: gwmux,
	}

	log.Println("Starting up webserver on :8080")
	log.Fatalln(gwServer.ListenAndServe())
}
