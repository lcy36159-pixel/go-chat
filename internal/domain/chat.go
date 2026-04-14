package domain

import "time"

type Chat struct {
	ID        uint
	Name      string
	Type      string
	CreatedBy uint
	CreatedAt time.Time
}

type ChatMember struct {
	ChatID uint
	UserID uint
}
