package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var listen = flag.String("listen", ":3333", "bind address")

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
