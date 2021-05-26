package models

import "gorm.io/gorm"

type User struct {
	Username string
	Password string
	Email    string

	gorm.Model
}

func (User) TableName() string {
	return "users"
}

type Wallet struct {
	Name	        string
	User            string

	gorm.Model
}

func (Wallet) TableName() string {
	return "wallet"
}
