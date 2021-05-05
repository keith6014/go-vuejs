package main

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	pb "go-vuejs/proto/helloworld"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

//go:embed "helloworld/helloworld.swagger.json"
var myjson string

//go:embed "helloworld/swagger-ui-3.47.1/dist"
var content embed.FS

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

	//	gwServer := &http.Server{
	//		Addr:    ":8080",
	//		Handler: gwmux,
	//	}

	//fs := http.FileServer(http.Dir("./helloworld"))
	//	err = gwmux.HandlePath("GET", "/b", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	//		w.Write([]byte("hello " + pathParams["name"]))
	//	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//	mux := http.NewServeMux()
	//	mux.Handle("/", gwmux)
	//	mux.Handle("/helloworld/", http.StripPrefix("/helloworld/", fs))
	//	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
	//		io.WriteString(w, string(myjson))
	//	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/", gwmux)
	r.HandleFunc("/api/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(myjson)
		if err != nil {
			http.Error(w, err.Error(), 422)
		}
		w.Write(b)
	})
	r.Handle("/d", clientHandler())

	log.Println("Starting up webserver on :8080")
	//log.Fatalln(gwServer.ListenAndServe())
	log.Println(http.ListenAndServe(":8080", r))
}

func clientHandler() http.Handler {
	fsys := fs.FS(content)
	contentStatic, err := fs.Sub(fsys, "index.html")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(contentStatic)
	return http.FileServer(http.FS(contentStatic))
}
