package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var listen = flag.String("listen", ":3333", "bind address")

//go:embed static/swagger-ui-3.47.1/dist
var static embed.FS

func main() {
	flag.Parse()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))

	})

	log.Println("Starting on", *listen)
	log.Fatal(http.ListenAndServe(*listen, r))
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		log.Println("FileServer does not permit any URL parameters.")
		//		panic("FileServer does not permit any URL parameters.")
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

type Status struct {
	Version string `json:"version,omitempty"`
}

func NewStatus() *Status {
	return &Status{Version: "0.0.0"}
}

func (s *Status) Render(w http.ResponseWriter, r *http.Request) error {
	s.Version = "1.1.1"
	return nil
}
