package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
	Email    string

	Transactions []Transaction
	Accounts     []Account
	Categories   []Category

	gorm.Model
}

func (User) TableName() string {
	return "users"
}

const (
	CATEGORIES_INCOME  = "income"
	CATEGORIES_EXPENSE = "expense"
)

type Category struct {
	Name            string
	TransactionType string
	UserID          uint
	User            User

	gorm.Model
}

func (Category) TableName() string {
	return "categories"
}

type Account struct {
	Name   string
	UserID uint
	User   User

	Transactions []Transaction

	gorm.Model
}

func (Account) TableName() string {
	return "accounts"
}

type Transaction struct {
	Date       time.Time
	AccountID  uint
	Account    Account
	CategoryID uint
	Category   Category
	UserID     uint
	User       User
	Amount     float64
	Note       string
	ImgURL     string

	gorm.Model
}

func (Transaction) TableName() string {
	return "transactions"
}
