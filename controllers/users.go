package controllers

import (
	"jiji/models"
	"jiji/utils"
	"jiji/views"
	"log"
	"net/http"
)

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
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

func NewUsers(us models.UserService) *Users {
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
	var vd views.Data
	var signupForm SignupForm

	err := parseForm(r, &signupForm)
	if err != nil {
		log.Println(err)
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: views.AlertMsgGeneric,
		}
		u.NewView.Render(w, vd)
		return
	}

	user := models.User{
		Username: signupForm.Username,
		Email:    signupForm.Email,
		Password: signupForm.Password,
	}

	err = u.us.Create(&user)
	if err != nil {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlError,
			Message: err.Error(),
		}
		u.NewView.Render(w, vd)
		return
	}

	err = u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
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

	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// Set token if the user doesn't have one and set it to cookie
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Token == "" {
		token, err := utils.GenerateToken()
		if err != nil {
			return err
		}
		user.Token = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "authToken",
		Value:    user.Token,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}
