package router

import (
	"net/http"

	"library-api/internal/handler"
)

// Router contiene la configuración de rutas HTTP.

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.HealthHandler)

	return mux
}
