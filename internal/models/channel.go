package models

import (
	"database/sql/driver"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ChannelType string

const (
	ChannelTypeOpen    ChannelType = "O" // ChannelTypeOpen
	ChannelTypePrivate ChannelType = "P" // ChannelTypePrivate
	ChannelTypeDirect  ChannelType = "D" // ChannelTypeDirect
	ChannelTypeGroup   ChannelType = "G" // ChannelTypeGroup
)

var ChannelTypesMap = map[string]ChannelType{
	"O": ChannelTypeOpen,
	"P": ChannelTypePrivate,
	"D": ChannelTypeDirect,
	"G": ChannelTypeGroup,
}

func (ct *ChannelType) Scan(value interface{}) error {
	*ct = ChannelType(value.([]byte))
	return nil
}

func (ct ChannelType) Value() (driver.Value, error) {
	return string(ct), nil
}

type Channel struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`

	Name          string      `json:"name"`
	Type          ChannelType `json:"type" sql:"type:ENUM('O', 'P', 'D', 'G')" gorm:"column:type"`
	TotalMsgCount uint64      `json:"totalMsgCount"`
}

func (c *Channel) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}

	return
}
