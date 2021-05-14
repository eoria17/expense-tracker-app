package controllers

import (
	"fmt"
	"net/http"

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
	r.HandleFunc("/register", ae.Register)
	r.HandleFunc("/home", ae.Home)
}

func (ae AppEngine) GetUser(username string) (CurrentUser *models.User) {
	db := ae.Storage.DB

	db.Where("username = ?", username).First(&CurrentUser)

	return
}

func (ae AppEngine) Home(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte("Hello, " + session.Values["username"].(string) + "!"))
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
