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
	r.POST("/chats/private", handler.CreatePrivateChatHandler)
	r.POST("/chats/group", handler.CreateGroupChatHandler)
	r.POST("/chats/:id/read", handler.MarkReadHandler)
	r.GET("/ws", handler.WebSocketHandler)
	r.GET("/messages", handler.GetMessagesHandler)
	r.GET("/chats", handler.GetChatsHandler)
	r.Run(":8080")
}
