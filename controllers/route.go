package controllers

import (
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/gorilla/mux"
)

type AppEngine struct {
	Storage *config.Storage
}

func (ae AppEngine) Route(r *mux.Router) {
	r.HandleFunc("/", ae.Login)
}
