package models

import (
	"gorm.io/gorm"
	"time"
)

type Posts struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID    string `gorm:"size:191"`
	User      User
	ChannelID string `gorm:"size:191"`
	Channel   Channel
}
