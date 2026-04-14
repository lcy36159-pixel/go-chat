package main

import (
	"go-chat/internal/handler"
	"go-chat/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	// db.DB.AutoMigrate(
	// 	&domain.User{},
	// 	&domain.Chat{},
	// 	&domain.ChatMember{},
	// 	&domain.Message{},
	// 	&domain.MessageRead{},
	// )

	r := gin.Default()
	r.POST("/chats/group", handler.CreateGroupChatHandler)
	r.GET("/ws", handler.WebSocketHandler)
	r.GET("/messages", handler.GetMessagesHandler)
	r.Run(":8080")
}
