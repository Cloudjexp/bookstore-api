package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status: "ok",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "error al generar la respuesta", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("servidor iniciado en http://localhost:8080")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
