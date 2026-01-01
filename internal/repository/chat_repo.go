package repository

import (
	"context"

	"github.com/BalaadityaPatanjali/ourDojo/internal/db"
)

func GetSingleConversationID(ctx context.Context) (string, error) {
	var id string

	err := db.Pool.QueryRow(
		ctx,
		`SELECT id FROM conversations LIMIT 1`,
	).Scan(&id)

	if err == nil {
		return id, nil
	}

	// If not found, create it ONCE
	err = db.Pool.QueryRow(
		ctx,
		`INSERT INTO conversations (id) VALUES (gen_random_uuid()) RETURNING id`,
	).Scan(&id)

	return id, err
}

