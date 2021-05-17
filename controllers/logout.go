package controllers

import (
	"fmt"
	"net/http"

	"github.com/eoria17/expense-tracker-app/models"
)

func (ae AppEngine) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	emptyUser := models.User{}

	session.Options.MaxAge = -1
	session.Values["logged_in"] = false
	session.Values["user"] = &emptyUser

	err = session.Save(r, w)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
