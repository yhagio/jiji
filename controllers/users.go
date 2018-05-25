package controllers

import (
	"fmt"
	"jiji/models"
	"jiji/views"
	"net/http"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

type SignupForm struct {
	Username string `schema:"username"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
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

// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var loginForm LoginForm
	if err := parseForm(r, &loginForm); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:  "Email",
		Value: user.Email,
	}
	http.SetCookie(w, &cookie)
	fmt.Fprint(w, "Login", user)
}
