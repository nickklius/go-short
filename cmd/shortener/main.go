package main

import "github.com/nickklius/go-short/internal/service"

func main() {
	s := service.NewService()

	s.Start()
}
