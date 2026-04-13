package domain

import "time"

type Message struct {
	ID        uint `gorm:"primaryKey"`
	SenderID  string
	TargetID  string
	ChatType  string
	Content   string
	CreatedAt time.Time
	Read      bool
}
