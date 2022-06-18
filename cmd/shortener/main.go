package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nickklius/go-short/internal/service"
)

func main() {
	errCh := run()

	if err := <-errCh; err != nil {
		log.Fatal(err)
	}
}

func run() <-chan error {
	errCh := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	s, err := service.NewService(ctx)
	if err != nil {
		errCh <- err
	}

	s.Start(ctx, errCh)

	go func() {
		<-ctx.Done()

		defer func() {
			stop()
			close(errCh)
		}()
	}()

	return errCh
}
