package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/eoria17/expense-tracker-app/config"
	"github.com/eoria17/expense-tracker-app/controllers"
	"github.com/eoria17/expense-tracker-app/models"
	"github.com/gorilla/mux"
)

func main() {
	//connect to DB
	storage_, err := config.Open(config.DB_HOST, config.DB_NAME, config.DB_USER, config.DB_PASSWORD, config.DB_PORT)
	if err != nil {
		panic(err)
	}

	//create AWS session
	creds := credentials.NewStaticCredentials(config.ACCESS_KEY_ID, config.SECRET_ACCESS_KEY, "")
	creds.Get()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AWS_REGION),
		Credentials: creds,
	})

	if err != nil {
		fmt.Println(err)
	}

	//S3 AWS client
	//s3c := s3.New(sess)
	s3uploader := s3manager.NewUploader(sess)


	//cognito client
	cognito := cognitoidentityprovider.New(sess)

	//dependency injection
	appEngine := controllers.AppEngine{
		Storage:  storage_,
		S3Client: s3uploader,
		Cognito: cognito,
	}

	//create route handler
	router := mux.NewRouter()
	appEngine.Route(router)

	//serve public as static file
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/"))))
	http.Handle("/assets/", router)

	//specify gob for sessions
	gob.Register(models.User{})

	//run server
	fmt.Println("Currently Listening to port 5000..")
	log.Println(http.ListenAndServe(":5000", router))
}
