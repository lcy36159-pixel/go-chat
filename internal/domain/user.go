package domain

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:50;uniqueIndex;not null"`
	Password  string    `gorm:"column:password;type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
