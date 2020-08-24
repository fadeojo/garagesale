package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	//Convert the Echo function to a type that implements http.Handler
	handler := http.HandlerFunc(Echo)

	if err := http.ListenAndServe("localhost:8004", handler); err != nil {
		log.Fatalf("error: listening and serving: %s \n", err)
	}

}

//Echo is a basic HTP handler
func Echo(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "You have asked to %s %s \n", r.Method, r.URL.Path); err != nil {
		 log.Fatalf("Error while writing response: %s \n", err)
	}
}
