package models

type Channel struct {
	ID    uint64 `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Users []User `gorm:"many2many:user_channels;"`
}
