package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	"github.com/gorilla/mux"
)

func (ae AppEngine) CreateExpenseTrx(w http.ResponseWriter, r *http.Request) {
	db := ae.Storage.DB

	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if not logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok && !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	viewPage := "views/createExpenseTrx.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	//variable
	date := time.Now()
	username := session.Values["user"].(models.User).Username
	accounts := []models.Account{}
	categories := []models.Category{}
	account_id := 0
	category_id := 0
	amount := 0.0
	error_message := ""
	error_message_bool := false
	user_id := session.Values["user"].(models.User).ID

	//get
	db.Where("user_id = ?", session.Values["user"].(models.User).ID).Find(&accounts)
	db.Where("user_id = ? AND transaction_type = ?", session.Values["user"].(models.User).ID, models.CATEGORIES_EXPENSE).Find(&categories)

	if r.Method == "POST" {

		//get image
		file, fileheader, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}

		//image uploader
		if file != nil {
			defer file.Close()

			size := fileheader.Size

			r.ParseMultipartForm(size)

			buffer := make([]byte, size)
			file.Read(buffer)

			fileBytes := bytes.NewReader(buffer)
			fileType := http.DetectContentType(buffer)

			if fileType != "image/jpeg" {
				error_message = "File uploaded is not an image. Please upload an image. ([filename].jpg)"
				error_message_bool = true
			}
			// } else if fileType != "image/png" {
			// 	error_message = "File uploaded is not an image. Please upload an image. ([filename].png)"
			// 	error_message_bool = true
			// }

			if !error_message_bool {
				_, err = ae.S3Client.Upload(&s3manager.UploadInput{
					Bucket: aws.String(config.AWS_BUCKET_NAME),
					Key:    aws.String(username + "_" + date.Format("2006-01-02 15:04:05") + ".jpg"),
					Body:   fileBytes,
					ACL:    aws.String("public-read"),
				})

				if err != nil {
					fmt.Println(err)
				}
			}
		}

		//get category id using category name and user_id
		categoryName := r.FormValue("category")
		if categoryName == "" {
			error_message = "Category must not be empty."
			error_message_bool = true
		}

		account_id, _ = strconv.Atoi(r.FormValue("account"))
		amount, err = strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil || amount <= 0 {
			error_message = "Amount must not be empty and must be greater than 0."
			error_message_bool = true
		}

		if !error_message_bool {
			ae.Storage.DB.Raw("SELECT id FROM categories WHERE name = ? and user_id = ?", categoryName, user_id).Scan(&category_id)

			//if category doesn't exist create it
			if category_id == 0 && categoryName != "" {
				//create category
				newCategory := models.Category{
					Name:            categoryName,
					TransactionType: "expense",
					UserID:          session.Values["user"].(models.User).ID,
				}

				db.Create(&newCategory)
				category_id = int(newCategory.ID)
			}

			//save transaction

			imgUrl := ""

			if file != nil {
				imgUrl = config.AWS_IMG_PATH + username + "_" + date.Format("2006-01-02") + "+" + date.Format("15") + "%3A" + date.Format("04") + "%3A" + date.Format("05") + ".jpg"
			}

			trx := models.Transaction{
				Date:       date,
				AccountID:  uint(account_id),
				CategoryID: uint(category_id),
				UserID:     session.Values["user"].(models.User).ID,
				Amount:     amount,
				Note:       r.FormValue("notes"),
				ImgURL:     imgUrl,
			}
			db.Create(&trx)

			//update account balance
			//get account from database
			account := models.Account{}
			ae.Storage.DB.Where("id = ?", account_id).First(&account)

			account.Amount = account.Amount - amount
			db.Save(&account)

			//send user back to home page
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

	}

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"assets":             assetsUrl,
		"title":              "Create Transaction",
		"username":           username,
		"date":               date,
		"Accounts":           accounts,
		"Categories":         categories,
		"error_message":      error_message,
		"error_message_bool": error_message_bool,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "create_expense_trx", data)
}

func (ae AppEngine) CreateIncomeTrx(w http.ResponseWriter, r *http.Request) {
	db := ae.Storage.DB

	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if not logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok && !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	viewPage := "views/createIncomeTrx.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	//variable
	date := time.Now()
	username := session.Values["user"].(models.User).Username
	accounts := []models.Account{}
	categories := []models.Category{}
	account_id := 0
	category_id := 0
	amount := 0.0
	error_message := ""
	error_message_bool := false
	user_id := session.Values["user"].(models.User).ID

	//get
	db.Where("user_id = ?", session.Values["user"].(models.User).ID).Find(&accounts)
	db.Where("user_id = ? AND transaction_type = ?", session.Values["user"].(models.User).ID, models.CATEGORIES_INCOME).Find(&categories)

	if r.Method == "POST" {

		//get image
		file, fileheader, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}

		//image uploader
		if file != nil {
			defer file.Close()

			size := fileheader.Size

			r.ParseMultipartForm(size)

			buffer := make([]byte, size)
			file.Read(buffer)

			fileBytes := bytes.NewReader(buffer)
			fileType := http.DetectContentType(buffer)

			if fileType != "image/jpeg" {
				error_message = "File uploaded is not an image. Please upload an image. ([filename].jpg)"
				error_message_bool = true
			}
			// } else if fileType != "image/png" {
			// 	error_message = "File uploaded is not an image. Please upload an image. ([filename].png)"
			// 	error_message_bool = true
			// }

			if !error_message_bool {
				_, err = ae.S3Client.Upload(&s3manager.UploadInput{
					Bucket: aws.String(config.AWS_BUCKET_NAME),
					Key:    aws.String(username + "_" + date.Format("2006-01-02 15:04:05") + ".jpg"),
					Body:   fileBytes,
					ACL:    aws.String("public-read"),
				})

				if err != nil {
					fmt.Println(err)
				}
			}
		}

		//get category id using category name and user_id
		categoryName := r.FormValue("category")
		if categoryName == "" {
			error_message = "Category must not be empty."
			error_message_bool = true
		}

		account_id, _ = strconv.Atoi(r.FormValue("account"))
		amount, err = strconv.ParseFloat(r.FormValue("amount"), 64)
		if err != nil || amount <= 0 {
			error_message = "Amount must not be empty and must be greater than 0."
			error_message_bool = true
		}

		if !error_message_bool {
			ae.Storage.DB.Raw("SELECT id FROM categories WHERE name = ? and user_id = ?", categoryName, user_id).Scan(&category_id)

			//if category doesn't exist create it
			if category_id == 0 && categoryName != "" {
				//create category
				newCategory := models.Category{
					Name:            categoryName,
					TransactionType: "income",
					UserID:          session.Values["user"].(models.User).ID,
				}

				db.Create(&newCategory)
				category_id = int(newCategory.ID)
			}

			//save transaction

			imgUrl := ""

			if file != nil {
				imgUrl = config.AWS_IMG_PATH + username + "_" + date.Format("2006-01-02") + "+" + date.Format("15") + "%3A" + date.Format("04") + "%3A" + date.Format("05") + ".jpg"
			}

			trx := models.Transaction{
				Date:       date,
				AccountID:  uint(account_id),
				CategoryID: uint(category_id),
				UserID:     session.Values["user"].(models.User).ID,
				Amount:     amount,
				Note:       r.FormValue("notes"),
				ImgURL:     imgUrl,
			}
			db.Create(&trx)

			//update account balance
			//get account from database
			account := models.Account{}
			ae.Storage.DB.Where("id = ?", account_id).First(&account)

			account.Amount = account.Amount + amount
			db.Save(&account)

			//send user back to home page
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

	}

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"assets":             assetsUrl,
		"title":              "Create Transaction",
		"username":           username,
		"date":               date,
		"Accounts":           accounts,
		"Categories":         categories,
		"error_message":      error_message,
		"error_message_bool": error_message_bool,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "create_income_trx", data)
}

func (ae AppEngine) Transaction(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if not logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok && !auth {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err != nil {
		fmt.Println(err)
	}

	//get transaction id from url
	vars := mux.Vars(r)
	transaction_id := vars["ID"]

	//get transaction data from database
	transactionData := models.Transaction{}
	ae.Storage.DB.Where("transactions.id = ?", transaction_id).First(&transactionData)

	//get transaction category
	category := models.Category{}
	ae.Storage.DB.Where("categories.id = ?", transactionData.CategoryID).First(&category)

	viewPage := "views/transaction.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	username := ""

	username = session.Values["user"].(models.User).Username
	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"title":           "Transaction",
		"assets":          assetsUrl,
		"username":        username,
		"transactionData": transactionData,
		"category":        category,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "transaction", data)
}
