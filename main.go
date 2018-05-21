package main

import (
	"jiji/views"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	contactView *views.View
	signupView  *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := signupView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("bootstrap",
		"views/home.gohtml")
	contactView = views.NewView("bootstrap",
		"views/contact.gohtml")
	signupView = views.NewView("bootstrap",
		"views/signup.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/signup", signup)
	http.ListenAndServe(":3000", r)
}
