package websocket

import (
	"net/http"
	"strconv"


	"github.com/gorilla/websocket"

	"github.com/BalaadityaPatanjali/ourDojo/internal/auth"
	"github.com/BalaadityaPatanjali/ourDojo/internal/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // tighten in prod
	},
}

type IncomingMessage struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	MediaURL string `json:"media_url"`
}

func ChatWS(w http.ResponseWriter, r *http.Request) {

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ParseToken(tokenString)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	userID := claims.UserID

	ChatHub.AddClient(userID, conn)
	defer func() {
		ChatHub.RemoveClient(userID)
		conn.Close()
	}()

	convID, err := repository.GetSingleConversationID(r.Context())
	if err != nil {
		http.Error(w, "conversation not found", http.StatusInternalServerError)
		return
	}

	// Send history
	history, err := repository.GetLastMessages(r.Context(), convID, 100)
	if err == nil {
		for _, msg := range history {
			msg["from_self"] = strconv.FormatBool(msg["sender_id"] == userID)
			conn.WriteJSON(msg)
		}
	}

	// Read loop
	for {
		var msg IncomingMessage
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		_ = repository.SaveMessage(
			r.Context(),
			convID,
			userID,
			msg.Type,
			msg.Content,
			msg.MediaURL,
		)

		ChatHub.SendToOther(userID, msg)
	}
}
