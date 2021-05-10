package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/controllers"
	"github.com/gorilla/mux"
)

func main() {
	//connect to DB
	storage_, err := config.Open(config.DB_HOST, config.DB_NAME, config.DB_USER, config.DB_PASSWORD, config.DB_PORT)
	if err != nil {
		panic(err)
	}

	//dependency injection
	appEngine := controllers.AppEngine{
		Storage: storage_,
	}

	//create route handler
	router := mux.NewRouter()
	appEngine.Route(router)

	//run server
	fmt.Println("Currently Listening to port 8080..")
	log.Println(http.ListenAndServe(":8080", router))
}
