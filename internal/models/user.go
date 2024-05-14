package models

type User struct {
	ID       uint64    `json:"id" gorm:"primary_key"`
	Name     string    `json:"name"`
	Email    string    `json:"email" gorm:"unique"`
	Password string    `json:"-"`
	Salt     string    `json:"-"`
	Channels []Channel `gorm:"many2many:user_channels;"`
	Contacts []User    `gorm:"many2many:user_contacts;"`
}
