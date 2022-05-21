package main

import (
	"log"

	"github.com/nickklius/go-short/internal/service"
)

func main() {
	s, err := service.NewService()
	if err != nil {
		log.Fatal(err)
	}

	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
