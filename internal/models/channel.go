package models

import (
	"database/sql/driver"
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
	BaseModel

	Name          string      `json:"name" gorm:"index:,class:FULLTEXT"`
	Type          ChannelType `json:"type" sql:"type:ENUM('O', 'P', 'D', 'G')" gorm:"column:type"`
	TotalMsgCount uint64      `json:"totalMsgCount"`
}
