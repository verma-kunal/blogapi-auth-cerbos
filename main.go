package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/verma-kunal/blogapi-auth-cerbos/service"
)

func main() {

	cerbosAddr := flag.String("cerbos", "localhost:3593", "Address of the Cerbos server")
	flag.Parse()

	// start the service API
	svc, err := service.New(*cerbosAddr)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}

	srv := http.Server{
		Addr:    ":8080",
		Handler: svc.Handler(),
	}

	log.Printf("Listening on %s", ":8080")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
