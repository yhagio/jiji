package controllers

import (
	"fmt"
	"jiji/views"
	"net/http"
)

type Users struct {
	NewView *views.View
}

type SignupForm struct {
	Username string `schema:"username"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
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
	var signupForm SignupForm
	if err := parseForm(r, &signupForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Username is", signupForm.Username)
	fmt.Fprintln(w, "Email is", signupForm.Email)
	fmt.Fprintln(w, "Password is", signupForm.Password)
}
