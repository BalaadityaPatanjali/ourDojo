package server

import (
	"log"
	"net/http"

	"github.com/BalaadityaPatanjali/ourDojo/internal/handlers"
)

func Start() {
	http.HandleFunc("/register", handlers.Register)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
