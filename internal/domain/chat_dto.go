package domain

import "time"

type ChatDTO struct {
	ChatID      uint
	Name        string
	LastMessage string
	UpdatedAt   time.Time
	UnreadCount int
}

// MemberInfo holds basic public information about a chat member.
type MemberInfo struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}
