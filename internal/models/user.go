package models

type User struct {
	BaseModel

	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Salt     string `json:"-"`
}
