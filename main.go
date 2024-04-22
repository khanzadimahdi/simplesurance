package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ringbuffer/counter"
	"ringbuffer/handlers"
	"sync"
	"syscall"
)

const (
	counterTTL      = 60
	counterFileName = "./storage/counter.json"
	httpServerPort  = 8080
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	c := counter.NewCounter(counterTTL)
	loadFromFile(c)
	defer storeToFile(c)

	var wg sync.WaitGroup
	server := startHttpServer(&wg, c)
	<-ctx.Done()

	server.Shutdown(ctx)
	wg.Wait()
}

func startHttpServer(wg *sync.WaitGroup, counter counter.Counter) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", handlers.NewCounterHandler(counter))

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", httpServerPort), Handler: mux}
	wg.Add(1)

	go func() {
		defer wg.Done()

		log.Println("server is up and running on port 8080")
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return server
}

func loadFromFile(c counter.Counter) {
	if _, err := os.Stat(counterFileName); errors.Is(err, os.ErrNotExist) {
		return
	}

	f, err := os.Open(counterFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := c.Load(f); err != nil {
		panic(err)
	}
}

func storeToFile(c counter.Counter) {
	f, err := os.Create(counterFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := c.Store(f); err != nil {
		panic(err)
	}
}
