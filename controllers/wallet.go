package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	"github.com/gorilla/mux"
)

func (ae AppEngine) Wallet(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if not logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok && !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	viewPage := "views/wallet.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	username := ""

	username = session.Values["user"].(models.User).Username

	//TODO
	vars := mux.Vars(r)
    wallet_id := vars["walletID"]

	//wallet_id := 1

	//get user's wallet from database
	wallet := models.Account{}
	ae.Storage.DB.Where("accounts.id = ?", wallet_id).First(&wallet)

	//fmt.Printf("%v", wallet)

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"title":       "Wallet",
		"assets":      assetsUrl,
		"username":    username,
		"walletData":  wallet,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "wallet", data)
}
