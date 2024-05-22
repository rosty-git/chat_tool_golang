package models

type Post struct {
	BaseModel

	UserID    string  `json:"user_id" gorm:"size:191"`
	User      User    `json:"user"`
	ChannelID string  `json:"channel_id" gorm:"size:191"`
	Channel   Channel `json:"-"`
	Message   string  `json:"message"`
}
