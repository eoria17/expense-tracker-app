package controllers

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	"github.com/jinzhu/now"
)

func (ae AppEngine) Home(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if not logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok && !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	db := ae.Storage.DB

	viewPage := "views/home.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	username := ""
	username = session.Values["user"].(models.User).Username
	user_id := session.Values["user"].(models.User).ID

	DashboardData := models.DashboardView{
		Transactions: []models.DashboardTransactionsView{},
	}
	transactions := []models.Transaction{}

	//fill in data
	DashboardData.MonthYear = time.Now().Format("January 2006")
	rows, err := db.Model(&models.Transaction{}).Where("user_id = ? AND date BETWEEN ? AND ?", user_id, now.BeginningOfMonth(), now.EndOfMonth()).Rows()
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		transaction := models.Transaction{}
		db.ScanRows(rows, &transaction)

		category := models.Category{}
		db.Where("id = ?", transaction.CategoryID).First(&category)

		account := models.Account{}
		db.Where("id = ?", transaction.AccountID).First(&account)

		if category.TransactionType == models.CATEGORIES_EXPENSE {
			DashboardData.Expenses += transaction.Amount
		} else if category.TransactionType == models.CATEGORIES_INCOME {
			DashboardData.Income += transaction.Amount
		}

		transaction.Category = category
		transaction.Account = account
		transactions = append(transactions, transaction)
	}

	//get time for transactions
	var dates []time.Time
	db.Raw("SELECT CAST(date AS DATE) FROM transactions GROUP BY CAST(date AS DATE) ORDER BY date DESC").Scan(&dates)

	for _, date := range dates {
		DashboardTransactionData := models.DashboardTransactionsView{
			Transactions: []models.Transaction{},
		}

		DashboardTransactionData.Day = date.Weekday().String()
		DashboardTransactionData.Date = date.Format("02-January-2006")

		totalIncome := 0.0
		totalExpenses := 0.0

		for _, transaction := range transactions {

			if date.Day() == transaction.Date.Day() {

				category := transaction.Category

				if category.TransactionType == models.CATEGORIES_EXPENSE {
					totalExpenses += transaction.Amount
				} else if category.TransactionType == models.CATEGORIES_INCOME {
					totalIncome += transaction.Amount
				}

				DashboardTransactionData.Transactions = append(DashboardTransactionData.Transactions, transaction)

				// transactions[i] = transactions[len(transactions)-1]
				// transactions[len(transactions)-1] = models.Transaction{}
				// transactions = transactions[:len(transactions)-1]

			}
		}

		DashboardTransactionData.TotalExpenses = totalExpenses
		DashboardTransactionData.TotalIncome = totalIncome
		DashboardData.Transactions = append(DashboardData.Transactions, DashboardTransactionData)
	}

	DashboardData.Total = DashboardData.Income - DashboardData.Expenses

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"title":         "Home",
		"assets":        assetsUrl,
		"username":      username,
		"DashboardData": DashboardData,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "home", data)
}
