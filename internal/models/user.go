package models

import "time"

type User struct {
	ID           string    // UUID from DB
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
