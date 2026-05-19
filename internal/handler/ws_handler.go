package handler

import (
	"encoding/json"
	"go-chat/internal/domain"
	"go-chat/internal/middleware"
	"go-chat/internal/repository"
	"go-chat/internal/service"
	"go-chat/internal/ws"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type IncomingMessage struct {
	ChatID  uint   `json:"chat_id"`
	Content string `json:"content"`
}

func WebSocketHandler(c *gin.Context) {
	userID, err := middleware.UserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// ✅ 升級 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		println("ReadMessage error:", err.Error())
		return
	}

	defer conn.Close()

	// ✅ 註冊 client
	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
	}
	ws.Register(client)

	for {
		// ✅ 讀取訊息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			println("ReadMessage error:", err.Error())
			break
		}

		println("收到:", string(msg))

		var incoming IncomingMessage
		if err := json.Unmarshal(msg, &incoming); err != nil {
			println("JSON error:", err.Error())
			continue
		}

		// ✅ 檢查聊天室成員
		userIDs, err := repository.GetUserIDsByChatID(incoming.ChatID)
		if err != nil {
			println("GetUserIDs error:", err.Error())
			continue
		}

		// 檢查是否為成員
		isMember := false
		for _, uid := range userIDs {
			if uid == userID {
				isMember = true
				break
			}
		}
		if !isMember {
			println("Not a member of chat:", incoming.ChatID)
			continue
		}

		// ✅ 建立 message
		message := domain.Message{
			ChatID:   incoming.ChatID,
			SenderID: &userID,
			Type:     "text",
			Content:  incoming.Content,
		}

		// ✅ 存 DB
		if err := service.HandleMessage(&message); err != nil {
			println("SaveMessage error:", err.Error())
			continue
		}

		// ✅ 序列化
		data, err := json.Marshal(message)
		if err != nil {
			println("Marshal error:", err.Error())
			continue
		}

		// ✅ 廣播給聊天室所有人
		for _, uid := range userIDs {
			ws.SendToUser(uid, data)
		}

		println("Broadcast success to", len(userIDs), "users")
	}
}
