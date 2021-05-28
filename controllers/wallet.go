package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
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

	//get user's wallets from database
	wallets := []models.Wallet{}
	ae.Storage.DB.Raw("SELECT name FROM wallet").Scan(&wallets)

	//code :=   route.Vars(r)["code"]



	//fmt.Printf("%v", wallets)

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"title":       "Wallet",
		"assets":      assetsUrl,
		"username":    username,
		"walletsData": wallets,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "wallets", data)
}
