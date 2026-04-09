package main

import (
	"go-chat/internal/domain"
	"go-chat/internal/handler"
	"go-chat/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()

	db.DB.AutoMigrate(&domain.Message{})

	r := gin.Default()

	r.GET("/ws", handler.WebSocketHandler)

	r.Run(":8080")
}