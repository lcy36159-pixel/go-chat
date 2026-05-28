package repository

import (
	"go-chat/internal/domain"
	"go-chat/pkg/db"
)

func CreateUser(user *domain.User) error {
	return db.DB.Create(user).Error
}

func GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
