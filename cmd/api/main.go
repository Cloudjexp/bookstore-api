package main

import (
	"log"
	"net/http"

	"library-api/internal/router"
)

func main() {
	mux := router.NewRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("servidor iniciado en http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
