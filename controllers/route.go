package controllers

import (
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type AppEngine struct {
	Storage  *config.Storage
	S3Client *s3manager.Uploader
}

func (ae AppEngine) Route(r *mux.Router) {
	r.HandleFunc("/", ae.Login)
	r.HandleFunc("/logout", ae.Logout)
	r.HandleFunc("/register", ae.Register)
	r.HandleFunc("/home", ae.Home)

	r.HandleFunc("/transaction/create/expense", ae.CreateExpenseTrx)
	r.HandleFunc("/transaction/create/income", ae.CreateIncomeTrx)
	r.HandleFunc("/wallets", ae.Wallets)
	r.HandleFunc("/wallet/create", ae.AddWallet)
	r.HandleFunc("/wallet/{walletID}", ae.Wallet)
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
