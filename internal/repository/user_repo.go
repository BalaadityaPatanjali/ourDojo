package repository

import (
	"context"
	"errors"

	"github.com/BalaadityaPatanjali/ourDojo/internal/db"
	"github.com/BalaadityaPatanjali/ourDojo/internal/models"
)

func CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	return db.Pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	user := &models.User{}

	err := db.Pool.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

