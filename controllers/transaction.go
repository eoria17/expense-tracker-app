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

	//get
	db.Where("user_id = ?", session.Values["user"].(models.User).ID).Find(&accounts)
	db.Where("user_id = ? AND transaction_type = ?", session.Values["user"].(models.User).ID, models.CATEGORIES_EXPENSE).Find(&categories)

	if r.Method == "POST" {

		//get image
		file, fileheader, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}
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
			//get category id using category name and user_id
			user_id := session.Values["user"].(models.User).ID
			categoryName := r.FormValue("category")

			category := -1
			//ae.Storage.DB.Select("categories.ID").Where("categories.name = ? and categories.user_id = ", categoryName, user_id).First(&category)
			ae.Storage.DB.Raw("SELECT id FROM categories WHERE name = ? and user_id = ?", categoryName, user_id).Scan(&category)

			//if category doesn't exist create it
			if category < 0{
				//create category
				newCategory := models.Category{
					Name: categoryName,
					TransactionType: "expense",
					UserID:     session.Values["user"].(models.User).ID,
				}
	
				db.Create(&newCategory)

				//get id of new category
				ae.Storage.DB.Raw("SELECT id FROM categories WHERE name = ?", categoryName).Scan(&category)
			}

			//save transaction

			account_id, _ = strconv.Atoi(r.FormValue("account"))
			category_id = category
			//category_id, _ = strconv.Atoi(r.FormValue("category"))
			amount, _ = strconv.ParseFloat(r.FormValue("amount"), 64)

			trx := models.Transaction{
				Date:       date,
				AccountID:  uint(account_id),
				CategoryID: uint(category_id),
				UserID:     session.Values["user"].(models.User).ID,
				Amount:     amount,
				Note:       r.FormValue("notes"),
				ImgURL:     config.AWS_IMG_PATH + username + "_" + date.Format("2006-01-02 15:04:05") + ".jpg",
			}

			fmt.Println(category_id)

			db.Create(&trx)
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

	//get
	db.Where("user_id = ?", session.Values["user"].(models.User).ID).Find(&accounts)
	db.Where("user_id = ? AND transaction_type = ?", session.Values["user"].(models.User).ID, models.CATEGORIES_INCOME).Find(&categories)

	if r.Method == "POST" {

		//get image
		file, fileheader, err := r.FormFile("image")
		if err != nil {
			fmt.Println(err)
		}
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
			//get category id using category name and user_id
			user_id := session.Values["user"].(models.User).ID
			categoryName := r.FormValue("category")

			category := -1
			//ae.Storage.DB.Select("categories.ID").Where("categories.name = ? and categories.user_id = ", categoryName, user_id).First(&category)
			ae.Storage.DB.Raw("SELECT id FROM categories WHERE name = ? and user_id = ?", categoryName, user_id).Scan(&category)

			//if category doesn't exist create it
			if category < 0{
				//create category
				newCategory := models.Category{
					Name: categoryName,
					TransactionType: "income",
					UserID:     session.Values["user"].(models.User).ID,
				}
	
				db.Create(&newCategory)

				//get id of new category
				ae.Storage.DB.Raw("SELECT id FROM categories WHERE name = ?", categoryName).Scan(&category)
			}


			//save transaction

			account_id, _ = strconv.Atoi(r.FormValue("account"))
			//category_id, _ = strconv.Atoi(r.FormValue("category"))
			category_id = category
			amount, _ = strconv.ParseFloat(r.FormValue("amount"), 64)

			trx := models.Transaction{
				Date:       date,
				AccountID:  uint(account_id),
				CategoryID: uint(category_id),
				UserID:     session.Values["user"].(models.User).ID,
				Amount:     amount,
				Note:       r.FormValue("notes"),
				ImgURL:     config.AWS_IMG_PATH + username + "_" + date.Format("2006-01-02 15:04:05") + ".jpg",
			}

			db.Create(&trx)
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

	if err != nil{
		fmt.Println(err)
	}

	


	viewPage := "views/transaction.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	username := ""

	username = session.Values["user"].(models.User).Username
	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"title":    "Transaction",
		"assets":   assetsUrl,
		"username": username,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "transaction", data)
}

