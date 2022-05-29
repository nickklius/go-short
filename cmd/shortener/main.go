package main

import (
	"context"
	"log"

	"github.com/nickklius/go-short/internal/service"
)

func main() {
	s, err := service.NewService(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
