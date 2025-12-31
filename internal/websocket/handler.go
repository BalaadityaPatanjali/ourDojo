package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/BalaadityaPatanjali/ourDojo/internal/auth"
	"github.com/BalaadityaPatanjali/ourDojo/internal/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // OK for dev
	},
}

type IncomingMessage struct {
	Type     string `json:"type"`       // text, emoji, image, etc.
	Content  string `json:"content"`    // text / emoji
	MediaURL string `json:"media_url"`  // Cloudinary URL
}

func ChatWS(w http.ResponseWriter, r *http.Request) {

	// 1. Read JWT from query
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	// 2. Validate JWT
	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// 3. Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	userID := claims.UserID

	// 4. Register client
	ChatHub.AddClient(userID, conn)
	defer func() {
		ChatHub.RemoveClient(userID)
		conn.Close()
	}()

	// 5. Read messages loop
	for {
		var msg IncomingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		// Always use the single conversation
		convID, err := repository.GetSingleConversationID(r.Context())
		if err != nil {
			continue
		}

		// Save message
		_ = repository.SaveMessage(
			r.Context(),
			convID,
			userID,
			msg.Type,
			msg.Content,
			msg.MediaURL,
		)

		// Send to partner
		ChatHub.SendToOther(userID, msg)
	}
}
