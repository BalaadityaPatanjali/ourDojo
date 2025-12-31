package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/BalaadityaPatanjali/ourDojo/internal/models"
	"github.com/BalaadityaPatanjali/ourDojo/internal/repository"
	"github.com/BalaadityaPatanjali/ourDojo/pkg/utils"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
	}

	// Store user
	err = repository.CreateUser(context.Background(), user)
	if err != nil {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
	})
}
