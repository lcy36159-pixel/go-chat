package handler

import (
	"go-chat/internal/domain"
	"go-chat/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(c *gin.Context) {
	// 把body的內容從JSON轉乘struct(主要讀取account、username和password)
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// 執行註冊邏輯
	userID, err := service.Register(req.Account, req.Username, req.Password)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 註冊成功回傳201(StatusCreated)和使用者ID
	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
	})
}

func LoginHandler(c *gin.Context) {
	// 讀取account和password
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// 執行登入邏輯
	token, err := service.Login(req.Account, req.Password)
	if err != nil {
		HandleError(c, err)
		return
	}
	// 登入成功回傳200(StatusOK)和JWT token
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
