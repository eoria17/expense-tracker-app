package controllers

import (
	"fmt"
	"net/http"

	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	"github.com/gorilla/mux"
)

type AppEngine struct {
	Storage *config.Storage
}

func (ae AppEngine) Route(r *mux.Router) {
	r.HandleFunc("/", ae.Login)
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
