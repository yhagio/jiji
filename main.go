package main

import (
	"fmt"
	"jiji/controllers"
	"jiji/models"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "jiji_dev_user"
	password = "123test"
	dbname   = "jiji_dev"
)

func main() {
	// Create a DB connection (Postgres)
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	// Create User Service
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticController.HomeView).Methods("GET")
	r.Handle("/contact", staticController.ContactView).Methods("GET")
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
