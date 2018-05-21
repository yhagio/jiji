package controllers

import (
	"fmt"
	"jiji/views"
	"net/http"

	"github.com/gorilla/schema"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Username string `schema:"username"`
	Email    string `schema:"email"`
	Passowrd string `schema:"password"`
}

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	decoder := schema.NewDecoder()
	signupForm := SignupForm{}
	if err := decoder.Decode(&signupForm, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprint(w, signupForm)
}
