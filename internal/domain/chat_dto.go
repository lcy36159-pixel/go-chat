package domain

import "time"

type ChatDTO struct {
	ChatID      uint
	Name        string
	LastMessage string
	UpdatedAt   time.Time
}
