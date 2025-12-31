package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BalaadityaPatanjali/ourDojo/internal/auth"
)

func Me(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r.Context())
	username := auth.GetUsername(r.Context())

	json.NewEncoder(w).Encode(map[string]string{
		"user_id":  userID,
		"username": username,
	})
}
