package main

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/webdevfuel/projectmotor/handler"
	"github.com/webdevfuel/projectmotor/template"
)

func main() {
	h := handler.NewHandler(handler.HandlerOptions{})
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fs := http.FileServer(http.Dir("./dist"))
	r.Handle("/dist/*", http.StripPrefix("/dist/", fs))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		template.Dashboard.ExecuteWriter(pongo2.Context{}, w)
	})
	r.Get("/login", h.Login)
	http.ListenAndServe("localhost:3000", r)
}
