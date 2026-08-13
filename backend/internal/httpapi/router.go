package httpapi

import (
	"log"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("POST /api/add", add)
	mux.HandleFunc("POST /api/subtract", subtract)
	mux.HandleFunc("POST /api/multiply", multiply)
	mux.HandleFunc("POST /api/divide", divide)
	mux.HandleFunc("POST /api/power", power)
	mux.HandleFunc("POST /api/square-root", squareRoot)
	mux.HandleFunc("POST /api/percentage", percentage)
	return mux
}

func health(w http.ResponseWriter, _ *http.Request) {
	if err := writeJSON(w, http.StatusOK, map[string]string{"status": "ok"}); err != nil {
		log.Printf("write health response: %v", err)
	}
}
