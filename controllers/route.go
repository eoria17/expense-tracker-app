package controllers

import (
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type AppEngine struct {
	Storage *config.Storage
}

func (ae AppEngine) Route(r *mux.Router) {
	r.HandleFunc("/", ae.Login)
	r.HandleFunc("/logout", ae.Logout)
	r.HandleFunc("/register", ae.Register)
	r.HandleFunc("/home", ae.Home)
	r.HandleFunc("/wallets", ae.Wallets)
	r.HandleFunc("/addWallet", ae.AddWallet)
	r.HandleFunc("/wallet", ae.Wallet)
}

func (ae AppEngine) GetUser(username string) (CurrentUser *models.User) {
	db := ae.Storage.DB

	db.Where("username = ?", username).First(&CurrentUser)

	return
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
