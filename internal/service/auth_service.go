package service

import (
	"errors"
	"fmt"
	"strings"

	"go-chat/internal/domain"
	"go-chat/internal/repository"
	"go-chat/pkg/jwt"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidRegisterInput = errors.New("account, username and password are required")
	ErrInvalidLoginInput    = errors.New("account and password are required")
	ErrWeakPassword         = errors.New("password must be at least 6 characters")
	ErrAccountTaken         = errors.New("account already exists")
	ErrInvalidCredentials   = errors.New("invalid account or password")
)

func Register(account, username, password string) (uint, error) {
	account = strings.TrimSpace(account)
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if account == "" || username == "" || password == "" {
		return 0, ErrInvalidRegisterInput
	}
	if len(password) < 6 {
		return 0, ErrWeakPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		Account:      account,
		Username:     username,
		PasswordHash: string(passwordHash),
	}
	if err := repository.CreateUser(user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrAccountTaken
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return user.ID, nil
}

func Login(account, password string) (string, error) {
	account = strings.TrimSpace(account)
	password = strings.TrimSpace(password)
	if account == "" || password == "" {
		return "", ErrInvalidLoginInput
	}

	user, err := repository.GetUserByAccount(account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
