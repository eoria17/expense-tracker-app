package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/gorilla/sessions"
)

var (
	key   = []byte(config.SESSION_KEY)
	store = sessions.NewCookieStore(key)
)

func (ae AppEngine) Login(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if logged in redirect to home
	if auth, ok := session.Values["logged_in"].(bool); ok && auth {
		http.Redirect(w, r, "/home", http.StatusFound)
	}

	viewPage := "views/login.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	login_err := ""
	password_err := ""
	username_err := ""
	username_err_bool := false
	password_err_bool := false
	login_err_bool := false
	username_filled := false
	username := ""

	if r.Method == "POST" {

		r.ParseForm()

		//check if username or password is null
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

		//search DB for login data
		if !username_err_bool && !password_err_bool {

			//login with cognito
			params := &cognitoidentityprovider.InitiateAuthInput{
				AuthFlow: aws.String("USER_PASSWORD_AUTH"),

				AuthParameters: map[string]*string{
					"USERNAME": aws.String(username),
					"PASSWORD": aws.String(r.FormValue("password")), //aws.String(r.FormValue("password")),
				},
				ClientId: aws.String(config.COGNITO_CLIENTID), // this is the app client ID
			}
			cog := ae.Cognito
			authResp, err := cog.InitiateAuth(params)

			fmt.Print(authResp)
			fmt.Print(err)

			user := ae.GetUser(r.FormValue("username"))

			if authResp == nil {
				login_err = "email or password is invalid"
				login_err_bool = true
			} else {
				session.Values["logged_in"] = true
				session.Values["user"] = user

				err = session.Save(r, w)
				if err != nil {
					fmt.Println(err)
				}

				//redirect to home
				http.Redirect(w, r, "/home", http.StatusFound)
				return
			}

		}

	}

	t, _ := template.ParseFiles(viewPage, config.HEADER_PATH, config.NAVIGATION_PATH)

	data := map[string]interface{}{
		"assets":            assetsUrl,
		"username_err_bool": username_err_bool,
		"password_err_bool": password_err_bool,
		"username_err":      username_err,
		"password_err":      password_err,
		"login_err":         login_err,
		"login_err_bool":    login_err_bool,
		"username_filled":   username_filled,
		"username":          username,
		"title":             "Login",
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "login", data)
}
