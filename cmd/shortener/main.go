package main

import (
	"github.com/nickklius/go-short/internal/app"
	"net/http"
)

func main() {

	http.HandleFunc("/", app.URLHandler)
	server := &http.Server{
		Addr: "localhost:8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
