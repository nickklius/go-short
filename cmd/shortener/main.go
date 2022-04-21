package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nickklius/go-short/internal/handlers"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	http.HandleFunc("/", handlers.URLHandler)
	http.ListenAndServe(":8080", r)
}
