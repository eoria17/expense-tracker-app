package main

import (
	"fmt"

	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	//connect to DB
	storage_, err := config.Open(config.DB_HOST, config.DB_NAME, config.DB_USER, config.DB_PASSWORD, config.DB_PORT)
	if err != nil {
		panic(err)
	}

	db := storage_.DB

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.Transaction{})

	DataSeed(db)
}

func DataSeed(db *gorm.DB) {
	//test user
	password, err := bcrypt.GenerateFromPassword([]byte("admin"), 14)
	if err != nil {
		fmt.Println(err)
	}

	user := models.User{
		Username: "admin",
		Email:    "admin@gmail.com",
		Password: string(password),
	}

	db.FirstOrCreate(&user)

	//default Categories (please add more later)
	category := models.Category{
		Name:            "Food",
		TransactionType: models.CATEGORIES_EXPENSE,
		UserID:          user.ID,
	}

	db.FirstOrCreate(&category)

	//test account
	account := models.Account{
		Name:   "Bank Account",
		UserID: user.ID,
	}

	db.FirstOrCreate(&account)

}
