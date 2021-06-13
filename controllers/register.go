package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/models"
	//"golang.org/x/crypto/bcrypt"
)

func (ae AppEngine) Register(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if logged in
	if auth, ok := session.Values["logged_in"].(bool); ok && auth {
		http.Redirect(w, r, "/main", http.StatusFound)
		return
	}

	//variable declarations
	viewPage := "views/register.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	username_err, email_err, password_err := "", "", ""
	username_err_bool, password_err_bool, email_err_bool := false, false, false
	username_filled, email_filled := false, false

	username, email := "", ""

	if r.Method == "POST" {
		r.ParseForm()

		//checks if empty
		if r.FormValue("email") == "" {
			email_err_bool = true
			email_err = "Please enter email."
		} else {
			email_filled = true
			email = r.FormValue("email")
		}

		if r.FormValue("username") == "" {
			username_err_bool = true
			username_err = "Please enter username."
		} else {
			username_filled = true
			username = r.FormValue("username")
		}

		if r.FormValue("password") == "" {
			password_err_bool = true
			password_err = "Please enter password."
		}

		//search DB for register data
		if !username_err_bool && !password_err_bool && !email_err_bool {
			user := ae.GetUser(r.FormValue("username"))

			tempUser := models.User{}
			ae.Storage.DB.Where("email = ?", email).First(&tempUser)

			if tempUser.Email != "" {
				email_err_bool = true
				email_err = "Email already exist, please enter a different email."
			}

			if user.Username != "" {
				username_err_bool = true
				username_err = "Username already exist, please enter a different username."
			} else {

				//password, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 14)
				if err != nil {
					fmt.Println(err)
				}

				newUser := models.User{
					Email:    r.FormValue("email"),
					Username: r.FormValue("username"),
					Password: string(""),
				}

				ae.Storage.DB.Create(&newUser)

				
			}

		}
		clientID := config.COGNITO_CLIENTID
		password := r.FormValue("password")

		signUpInput := &cognitoidentityprovider.SignUpInput{
			AnalyticsMetadata: &cognitoidentityprovider.AnalyticsMetadataType{},
			ClientId:          &clientID,
			ClientMetadata:    map[string]*string{},
			Password:          &password,
			UserAttributes:    []*cognitoidentityprovider.AttributeType{},
			UserContextData:   &cognitoidentityprovider.UserContextDataType{},
			Username:          &username,
			ValidationData:    []*cognitoidentityprovider.AttributeType{},
		}

		//create user in cognito
		ae.Cognito.SignUp(signUpInput)

		//redirect to login
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH)

	data := map[string]interface{}{
		"assets":            assetsUrl,
		"username_err_bool": username_err_bool,
		"password_err_bool": password_err_bool,
		"email_err_bool":    email_err_bool,
		"email_err":         email_err,
		"username_err":      username_err,
		"password_err":      password_err,
		"username_filled":   username_filled,
		"username":          username,
		"email_filled":      email_filled,
		"email":             email,
		"title":             "Register",
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "register", data)
}
