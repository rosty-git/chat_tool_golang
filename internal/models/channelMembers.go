package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ChannelMembers struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ChannelID    string `gorm:"size:191"`
	Channel      Channel
	UserID       string `gorm:"size:191"`
	User         User
	Roles        string
	LastViewedAt time.Time
	MsgCount     uint64
	MentionCount uint64
}

func (cm *ChannelMembers) BeforeCreate(tx *gorm.DB) (err error) {
	if cm.ID == "" {
		cm.ID = uuid.New().String()
	}

	return
}
