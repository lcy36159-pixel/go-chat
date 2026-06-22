package handler

import (
	"encoding/json"
	"go-chat/internal/service"
	"go-chat/internal/ws"
	"go-chat/pkg/jwt"
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
	// 升級 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		println("Upgrade error:", err.Error())
		return
	}
	// 確保連線在函式結束時關閉
	defer conn.Close()
	// 用 token 取得 userID
	token := c.Query("token")
	userID, err := jwt.ParseToken(token)
	if err != nil {
		println("Invalid token")
		return
	}

	// 註冊 client
	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
	}
	ws.Register(client)

	for {
		// 讀取訊息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			println("ReadMessage error:", err.Error())
			break
		}
		// 解析 JSON
		var incoming IncomingMessage
		if err := json.Unmarshal(msg, &incoming); err != nil {
			println("JSON error:", err.Error())
			continue
		}
		// 驗證、存檔並廣播
		if err := service.SendMessage(userID, incoming.ChatID, incoming.Content); err != nil {
			println("SendMessage error:", err.Error())
			continue
		}
	}
}
