package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/BalaadityaPatanjali/ourDojo/internal/auth"
	"github.com/BalaadityaPatanjali/ourDojo/internal/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // tighten later if needed
	},
}

type IncomingMessage struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	MediaURL string `json:"media_url"`
}

func ChatWS(w http.ResponseWriter, r *http.Request) {

	// 1. Read JWT
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

	// 5. Get single conversation
	convID, err := repository.GetSingleConversationID(r.Context())
	if err != nil {
		http.Error(w, "conversation not found", http.StatusInternalServerError)
		return
	}

	// 6. Send chat history
history, err := repository.GetLastMessages(r.Context(), convID, 100)
if err == nil {
	for _, msg := range history {

		// sender_id and userID are strings
		if msg["sender_id"] == userID {
			msg["from_self"] = "true"
		} else {
			msg["from_self"] = "false"
		}

		_ = conn.WriteJSON(msg)
	}
}

	// 7. Read loop
	for {
		var msg IncomingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
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

		// Send to other user
		ChatHub.SendToOther(userID, msg)
	}
}
