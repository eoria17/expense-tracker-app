package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
)

func (ae AppEngine) AddWallet(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if not logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok && !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	//check for submitted data
	if r.FormValue("name") != "" {
		//add new wallet to database
		newWallet := models.Wallet{
			Name:    r.FormValue("name"),
			User: session.Values["user"].(models.User).Username,
		}

		ae.Storage.DB.Create(&newWallet)

		//redirect to wallets page
		http.Redirect(w, r, "/wallets", http.StatusFound)

		return
	}


	viewPage := "views/addWallet.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	username := ""

	username = session.Values["user"].(models.User).Username

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"title":    "addWallet",
		"assets":   assetsUrl,
		"username": username,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "addWallet", data)
}
