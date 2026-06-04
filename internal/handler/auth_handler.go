package handler

import (
	"errors"
	"go-chat/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func RegisterHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, err := service.Register(req.Account, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRegisterInput), errors.Is(err, service.ErrWeakPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrAccountTaken):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed, please try again later"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
	})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	token, err := service.Login(req.Account, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidLoginInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed, please try again later"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
