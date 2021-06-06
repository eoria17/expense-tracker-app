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

func (ae AppEngine) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
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

	//get transaction id from url
	vars := mux.Vars(r)
    transaction_id := vars["transactionID"]

	//fmt.Println(wallet_id)

	//delete transaction
	ae.Storage.DB.Delete(&models.Transaction{}, transaction_id)

	//send user to wallets page
	http.Redirect(w, r, "/wallets", http.StatusFound)

}
