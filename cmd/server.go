package main

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"

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

type WebServer struct {
	r     chi.Router
	gwmux http.Handler
}

func NewWebServer(gatewayHandler http.Handler) *WebServer {
	srv := &WebServer{
		r:     chi.NewRouter(),
		gwmux: gatewayHandler,
	}
	srv.setup()
	return srv
}

func (srv *WebServer) setup() {
	r := srv.r

	fsys := fs.FS(static)
	swaggerStatic, err := fs.Sub(fsys, "static/swagger-ui-3.47.1/dist")
	if err != nil {
		log.Println("content static swagger", err)
	}

	beaufyStatic, err := fs.Sub(fsys, "static/beufy")
	if err != nil {
		log.Fatal("content static beaufy", err)
	}
	//	FileServer(r, "/api", http.FS(contentStatic))

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	//http://localhost:8080/static/buefy/
	FileServer(r, "/static/buefy", http.FS(beaufyStatic))

	FileServer(r, "/api", http.FS(swaggerStatic))

	r.HandleFunc("/api/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(swagger)
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("").Parse(indexTmpl)

		if err != nil {
			log.Println(err)
		}

		t.Execute(w, nil)

	})

	r.Handle("/hello_world", srv.gwmux)
}

func (srv *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.r.ServeHTTP(w, r)
}

//go:embed "helloworld/helloworld.swagger.json"
var swagger []byte

//go:embed static
var static embed.FS

//go:embed "tmpl/index.html"
var indexTmpl string

//go:embed node_modules
var node_modules embed.FS

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

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r := NewWebServer(gwmux)
	log.Println(http.ListenAndServe(":8080", r))
}

func clientHandler() http.Handler {
	fsys := fs.FS(static)
	contentStatic, err := fs.Sub(fsys, "index.html")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(contentStatic)
	return http.FileServer(http.FS(contentStatic))
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
