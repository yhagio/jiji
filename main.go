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
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery)

	r := mux.NewRouter()
	r.Handle("/", staticController.HomeView).Methods("GET")
	r.Handle("/contact", staticController.ContactView).Methods("GET")
	// Users
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")
	r.Handle("/login", usersController.LoginView).Methods("GET")
	r.HandleFunc("/login", usersController.Login).Methods("POST")
	// Galleries
	r.Handle("/galleries/new", galleriesController.New).Methods("GET")
	http.ListenAndServe(":3000", r)
}
