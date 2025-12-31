package repository

import (
	"context"

	"github.com/BalaadityaPatanjali/ourDojo/internal/db"
)

func GetSingleConversationID(ctx context.Context) (string, error) {
	var id string
	err := db.Pool.QueryRow(
		ctx,
		"SELECT id FROM conversations LIMIT 1",
	).Scan(&id)
	return id, err
}
