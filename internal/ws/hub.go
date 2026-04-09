package ws

import "github.com/gorilla/websocket"

type Client struct {
	UserID string
	Conn   *websocket.Conn
}

var Clients = make(map[string]*Client)

func Register(c *Client) {
	Clients[c.UserID] = c
}

func Broadcast(msg []byte) {
	for _, c := range Clients {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// ✅ 新增：發送給指定 user
func SendToUser(userID string, msg []byte) {
	if client, ok := Clients[userID]; ok {
		client.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
