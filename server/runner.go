package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Runner struct {
	server Server
}

func NewRunner(server Server) *Runner {
	return &Runner{server: server}
}

func (r *Runner) Run() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func(s Server) {
		defer wg.Done()
		if err := s.Start(); err != nil {
			log.Printf("server exited with error: %v", err)
		}
	}(r.server)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("received signal: %s, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.server.Stop(ctx); err != nil {
		log.Printf("error stopping server: %v", err)
	}

	wg.Wait()
	log.Println("all servers stopped gracefully")
}
