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
	user_id := session.Values["user"].(models.User).ID

	//get wallet ID from url
	vars := mux.Vars(r)
    wallet_id := vars["walletID"]

	//get user's wallet from database
	wallet := models.Account{}
	ae.Storage.DB.Where("accounts.id = ? AND accounts.user_id = ?", wallet_id, user_id).First(&wallet)

	//get wallet's transactions from database
	transactions := []models.Transaction{}
	ae.Storage.DB.Where("transactions.account_id = ? AND transactions.user_id = ?", wallet_id, user_id).Find(&transactions)

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"title":       "Wallet",
		"assets":      assetsUrl,
		"username":    username,
		"walletData":  wallet,
		"transactions": transactions,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "wallet", data)
}
