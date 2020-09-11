package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// =========================================================================
	// App Starting
	log.Printf("main: Started")
	defer log.Println("main: Completed")

	// =========================================================================
	// Start API Service
	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverError := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverError <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverError:
		log.Printf("error: listening and serving: %s", err)
		return

	case <-shutdown:
		log.Println("main: Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)

		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Printf("main : could not stop server gracefully : %v", err)
			return
		}
	}
}

// Echo is a basic HTTP Handler.
// If you open localhost:8000 in your browser, you may notice
// double request being made. This happens because the browser
// sends a request in the background for a website favicon.
func Echo(writer http.ResponseWriter, request *http.Request) {
	// Print a random number at the beginning and end of each request.
	n, err := rand.Int(rand.Reader, big.NewInt(1000))

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	log.Println("start", n)
	defer log.Println("end", n)

	// Simulate a long-running request.
	time.Sleep(3 * time.Second)

	_, _ = fmt.Fprintf(writer, "You asked to %s %s\n", request.Method, request.URL.Path)
}
