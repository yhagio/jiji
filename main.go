package main

import (
	"fmt"
	"jiji/controllers"
	"jiji/middlewares"
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
	requireUserMw := middlewares.RequireUser{
		UserService: services.User,
	}

	staticCtrl := controllers.NewStatic()
	usersCtrl := controllers.NewUsers(services.User)
	galleriesCtrl := controllers.NewGalleries(services.Gallery)

	r := mux.NewRouter()
	r.Handle("/", staticCtrl.HomeView).Methods("GET")
	r.Handle("/contact", staticCtrl.ContactView).Methods("GET")

	// Users
	r.HandleFunc("/signup", usersCtrl.New).Methods("GET")
	r.HandleFunc("/signup", usersCtrl.Create).Methods("POST")
	r.Handle("/login", usersCtrl.LoginView).Methods("GET")
	r.HandleFunc("/login", usersCtrl.Login).Methods("POST")

	// Galleries
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesCtrl.New)).Methods("GET")
	r.Handle("/galleries", requireUserMw.ApplyFunc(galleriesCtrl.Create)).Methods("POST")

	http.ListenAndServe(":3000", r)
}
