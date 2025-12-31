package repository

import (
	"context"

	"github.com/BalaadityaPatanjali/ourDojo/internal/db"
)

// SaveMessage persists a chat message to DB
func SaveMessage(
	ctx context.Context,
	conversationID string,
	senderID string,
	msgType string,
	content string,
	mediaURL string,
) error {

	_, err := db.Pool.Exec(
		ctx,
		`INSERT INTO messages (conversation_id, sender_id, type, content, media_url)
		 VALUES ($1, $2, $3, $4, $5)`,
		conversationID,
		senderID,
		msgType,
		content,
		mediaURL,
	)

	return err
}

func GetLastMessages(ctx context.Context, conversationID string, limit int) ([]map[string]string, error) {
	rows, err := db.Pool.Query(
		ctx,
		`SELECT sender_id, type, content, media_url, created_at
		 FROM messages
		 WHERE conversation_id = $1
		 ORDER BY created_at ASC
		 LIMIT $2`,
		conversationID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]string
	for rows.Next() {
		var senderID, msgType, content, mediaURL string
		var createdAt string

		_ = rows.Scan(&senderID, &msgType, &content, &mediaURL, &createdAt)

		messages = append(messages, map[string]string{
			"sender_id": senderID,
			"type":      msgType,
			"content":   content,
			"media_url": mediaURL,
		})
	}
	return messages, nil
}

