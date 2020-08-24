package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

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
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on: %s \n", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <- serverErrors:
		log.Fatalf("error listening and serving: %s", err)
	case <- shutdown:
		log.Println("main: Shutdown")

		// Give outstanding requests a deadline for completion.
		timeout := 5 * time.Second
		ctx,cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main: Graceful shutdown down did not happen in %v: %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main: could not stop server gracefully %v", err)
		}

	}

}

//Echo is a basic HTP handler
func Echo(w http.ResponseWriter, r *http.Request) {
	// Print a random number at the beginning and end of each request.
	n := rand.Intn(1000)
	log.Println("start", n)
	defer log.Println("stop", n)

	// Simulate a long-running request.
	time.Sleep(3 * time.Second)

	if _, err := fmt.Fprintf(w, "You have asked to %s %s \n", r.Method, r.URL.Path); err != nil {
		log.Fatalf("Error while writing response: %s \n", err)
	}
}
