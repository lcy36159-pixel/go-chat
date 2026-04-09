package service

import (
	"go-chat/internal/domain"
	"go-chat/internal/repository"
)

func HandleMessage(msg *domain.Message) error {
	return repository.SaveMessage(msg)
}