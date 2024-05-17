package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Salt     string `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	slog.Info("BeforeCreate", "u", u)

	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	return
}
