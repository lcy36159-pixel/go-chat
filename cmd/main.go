package main

import (
	"go-chat/internal/handler"
	"go-chat/internal/middleware"
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
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/chats/private", handler.CreatePrivateChatHandler)
	auth.POST("/chats/group", handler.CreateGroupChatHandler)
	auth.POST("/chats/:id/read", handler.MarkReadHandler)
	auth.GET("/ws", handler.WebSocketHandler)
	auth.GET("/messages", handler.GetMessagesHandler)
	auth.GET("/chats", handler.GetChatsHandler)
	r.Run(":8080")
}
