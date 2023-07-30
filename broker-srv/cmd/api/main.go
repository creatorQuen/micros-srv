package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8000"

type Config struct{}

func main() {
	application := Config{}

	log.Printf("Broker start %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: application.routes(),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
