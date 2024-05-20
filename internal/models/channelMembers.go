package models

import (
	"time"
)

type ChannelMembers struct {
	BaseModel

	ChannelID    string    `json:"channel_id" gorm:"size:191"`
	Channel      Channel   `json:"-"`
	UserID       string    `json:"user_id" gorm:"size:191"`
	User         User      `json:"-"`
	Roles        string    `json:"roles"`
	LastViewedAt time.Time `json:"last_viewed_at"`
	MsgCount     uint64    `json:"msg_count"`
	MentionCount uint64    `json:"mention_count"`
}
