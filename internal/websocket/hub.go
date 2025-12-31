package websocket

import "github.com/gorilla/websocket"

type Hub struct {
	Clients map[string]*websocket.Conn
}

var ChatHub = &Hub{
	Clients: make(map[string]*websocket.Conn),
}

func (h *Hub) AddClient(userID string, conn *websocket.Conn) {
	h.Clients[userID] = conn
}

func (h *Hub) RemoveClient(userID string) {
	delete(h.Clients, userID)
}

func (h *Hub) SendToOther(senderID string, payload any) {
	for uid, conn := range h.Clients {
		if uid != senderID {
			_ = conn.WriteJSON(payload)
		}
	}
}
