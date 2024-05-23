package models

import (
	"time"
)

type Status struct {
	BaseModel
	//CreatedAt time.Time      `json:"created_at"`
	//UpdatedAt time.Time      `json:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	UserID         string    `json:"user_id" gorm:"primaryKey"`
	Status         string    `json:"status"`
	Manual         bool      `json:"manual"`
	LastActivityAt time.Time `json:"last_activity_at"`
	DNDEndTime     time.Time `json:"dnd_end_time"`
	PrevStatus     string    `json:"prev_status"`
}
