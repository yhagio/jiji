package controllers

import (
	"fmt"
	"jiji/models"
	"jiji/views"
	"net/http"
)

type Users struct {
	NewView *views.View
	us      *models.UserService
}

type SignupForm struct {
	Username string `schema:"username"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
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
	user := models.User{
		Username: signupForm.Username,
		Email:    signupForm.Email,
		Password: signupForm.Password,
	}
	err := u.us.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, user)
}
