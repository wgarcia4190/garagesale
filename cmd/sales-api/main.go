package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wgarcia4190/garagesale/cmd/sales-api/internal/handlers"
	"github.com/wgarcia4190/garagesale/internal/platform/conf"
	"github.com/wgarcia4190/garagesale/internal/platform/database"
	"github.com/wgarcia4190/garagesale/internal/schema"
)

func main() {
	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error: generating config usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config %s", err)
	}

	// =========================================================================
	// App Starting
	log.Printf("main: Started")
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatalf("error : generating config for output : %v", err)
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Setup Dependencies
	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("applying migrations", err)
		}
		log.Println("Migrations complete")
		return

	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("applying seed data", err)
		}
		log.Println("Seed data inserted")
		return
	}

	// =========================================================================
	// Start API Service
	ps := handlers.Product{DB: db}
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(ps.GetListProducts),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
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
		timeout := cfg.Web.ShutdownTimeout
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
