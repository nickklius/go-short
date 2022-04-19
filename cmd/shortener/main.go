package main

import (
	"github.com/nickklius/go-short/internal/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/", handlers.URLHandler)
	server := &http.Server{
		Addr: "localhost:8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
