package server

import (
	"log"
	"net/http"

	"github.com/BalaadityaPatanjali/ourDojo/internal/auth"
	"github.com/BalaadityaPatanjali/ourDojo/internal/handlers"
	"github.com/BalaadityaPatanjali/ourDojo/internal/websocket"

)

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/ws", websocket.ChatWS)

	protected := http.HandlerFunc(handlers.Me)
	mux.Handle("/me", auth.JWTMiddleware(protected))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}

