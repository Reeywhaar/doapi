package main

import (
	"doctl/api"
	"doctl/cli"
	"doctl/internal"
	"log"
	"net/http"
	"os"
)

func main() {
	token := os.Getenv("DO_TOKEN")
	if token == "" {
		log.Fatal("DO_TOKEN environment variable is required")
	}

	client := internal.NewClient(token)

	if len(os.Args) > 1 {
		cli.Run(os.Args[1:], client)
		return
	}

	mux := http.NewServeMux()
	api.NewHandler(client).RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
