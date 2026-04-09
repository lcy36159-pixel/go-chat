package handler

import (
	"encoding/json"
	"go-chat/internal/domain"
	"go-chat/internal/service"
	"go-chat/internal/ws"
	"go-chat/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type IncomingMessage struct {
	SenderID string `json:"sender_id"`
	TargetID string `json:"target_id"`
	ChatType string `json:"chat_type"`
	Content  string `json:"content"`
}

func WebSocketHandler(c *gin.Context) {
	conn, _ := upgrader.Upgrade(c.Writer, c.Request, nil)

	userID := c.Query("user_id")

	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
	}

	ws.Register(client)

	// ✅ 補發離線訊息
	msgs, _ := repository.GetUnreadMessages(userID)

	for _, m := range msgs {
		data, _ := json.Marshal(m)
		ws.SendToUser(userID, data)
	}

	// 標記為已讀
	repository.MarkAsRead(userID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var incoming IncomingMessage
		json.Unmarshal(msg, &incoming)

		message := domain.Message{
			SenderID: incoming.SenderID,
			TargetID: incoming.TargetID,
			ChatType: incoming.ChatType,
			Content:  incoming.Content,
		}

		service.HandleMessage(&message)

		// ws.Broadcast(msg)
		ws.SendToUser(incoming.TargetID, msg)
	}
}
