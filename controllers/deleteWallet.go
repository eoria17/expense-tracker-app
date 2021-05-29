package controllers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	//"text/template"
	//"strconv"

	//"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
)

func (ae AppEngine) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if not logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok && !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err != nil{
		fmt.Println(err)
	}

	//get wallet id from url
	vars := mux.Vars(r)
    wallet_id := vars["walletID"]

	//fmt.Println(wallet_id)

	//delete wallet
	ae.Storage.DB.Delete(&models.Account{}, wallet_id)

	//TODO
	//delete transactions from wallet??

	//send user to wallets page
	http.Redirect(w, r, "/wallets", http.StatusFound)

}
