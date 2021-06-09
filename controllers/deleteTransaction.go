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
	//userID := session.Values["user"].(models.User).ID

	//get transaction id from url
	vars := mux.Vars(r)
    transaction_id := vars["transactionID"]

	user_id := session.Values["user"].(models.User).ID

	//get transaction from database
	transaction := models.Transaction{}
	ae.Storage.DB.Where("transactions.id = ? AND transactions.user_id = ?", transaction_id, user_id).First(&transaction)

	//get account from database
	account := models.Account{}
	ae.Storage.DB.Where("accounts.id = ?", transaction.AccountID).First(&account)

	//get transaction category from database
	category := models.Category{}
	ae.Storage.DB.Where("categories.id = ?", transaction.CategoryID).First(&category)

	//println("########### transaction type: " + category.TransactionType + "###############")

	//update account balance
	if category.TransactionType == "expense" {
		//do reverse action, we are 'undoing' the change because we are deleting transaction
		account.Amount = account.Amount + transaction.Amount
	} else if category.TransactionType == "income" {
		//do reverse action, we are 'undoing' the change because we are deleting transaction
		account.Amount = account.Amount - transaction.Amount
	}

	//save account change to db
	ae.Storage.DB.Save(&account)

	//delete transaction
	ae.Storage.DB.Delete(&models.Transaction{}, transaction_id)

	//send user to wallets page
	http.Redirect(w, r, "/wallets", http.StatusFound)

}
