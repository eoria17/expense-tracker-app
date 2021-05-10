package main

import (
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
)

func main() {
	//connect to DB
	storage_, err := config.Open(config.DB_HOST, config.DB_NAME, config.DB_USER, config.DB_PASSWORD, config.DB_PORT)
	if err != nil {
		panic(err)
	}

	db := storage_.DB

	db.AutoMigrate(&models.Test{})
}
