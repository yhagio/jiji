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

	r := mux.NewRouter()

	staticCtrl := controllers.NewStatic()
	usersCtrl := controllers.NewUsers(services.User)
	galleriesCtrl := controllers.NewGalleries(services.Gallery, r)

	requireUserMw := middlewares.RequireUser{
		UserService: services.User,
	}

	r.Handle("/", staticCtrl.HomeView).Methods("GET")
	r.Handle("/contact", staticCtrl.ContactView).Methods("GET")

	// Users
	r.HandleFunc("/signup", usersCtrl.New).Methods("GET")
	r.HandleFunc("/signup", usersCtrl.Create).Methods("POST")
	r.Handle("/login", usersCtrl.LoginView).Methods("GET")
	r.HandleFunc("/login", usersCtrl.Login).Methods("POST")

	// Galleries
	r.Handle("/galleries", requireUserMw.ApplyFunc(galleriesCtrl.GetAllByUser)).Methods("GET")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesCtrl.New)).Methods("GET")
	r.Handle("/galleries", requireUserMw.ApplyFunc(galleriesCtrl.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesCtrl.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFunc(galleriesCtrl.Edit)).Methods("GET")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFunc(galleriesCtrl.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFunc(galleriesCtrl.Delete)).Methods("POST")
	http.ListenAndServe(":3000", r)
}
