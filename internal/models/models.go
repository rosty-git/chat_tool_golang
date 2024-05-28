package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type BaseModel struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (bm *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if bm.ID == "" {
		bm.ID = uuid.New().String()
	}

	return
}
