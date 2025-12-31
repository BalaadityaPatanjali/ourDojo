package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/BalaadityaPatanjali/ourDojo/internal/db"
	"github.com/BalaadityaPatanjali/ourDojo/internal/server"
)

func main() {
	_ = godotenv.Load()

	if err := db.Connect(); err != nil {
		log.Fatal("DB connection failed:", err)
	}

	server.Start()
}
