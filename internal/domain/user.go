package domain

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Account      string    `gorm:"size:50;uniqueIndex;not null"`
	Username     string    `gorm:"size:50;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
