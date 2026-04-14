package domain

import "time"

type Message struct {
	ID        uint       `gorm:"primaryKey"`
	ChatID    uint       `gorm:"not null;index"` // 查詢會用到
	SenderID  *uint      `gorm:"index"`          // 用 pointer（對應 SET NULL）
	Type      string     `gorm:"type:message_type;default:'text'"`
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	DeletedAt *time.Time `gorm:"index"` // 軟刪除（收回訊息）
}
