package models

import "gorm.io/gorm"

type User struct {
	Name     string
	Username string
	Password string

	gorm.Model
}

func (User) TableName() string {
	return "user"
}
