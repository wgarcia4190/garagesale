package main

import (
	"context"
	"encoding/json"
	"log"
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
		Handler:      http.HandlerFunc(GetListProducts),
		ReadTimeout:  5 * time.Second, //nolint:gomnd
		WriteTimeout: 5 * time.Second, //nolint:gomnd
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

// Product is something we sell.
type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// GetListProducts gives all products as list.
func GetListProducts(writer http.ResponseWriter, request *http.Request) {
	list := []Product{
		{Name: "Comic Books", Cost: 75, Quantity: 50},
		{Name: "MCDonald's Toys", Cost: 25, Quantity: 120},
	}

	data, err := json.Marshal(list)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling", err)

		return
	}

	writer.Header().Set("content-type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(data); err != nil {
		log.Println("Error writing", err)
	}
}
