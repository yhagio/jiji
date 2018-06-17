package controllers

import (
	"jiji/email"
	"jiji/middlewares"
	"jiji/models"
	"jiji/utils"
	"jiji/views"
	"net/http"
	"time"
)

type Users struct {
	NewView            *views.View
	LoginView          *views.View
	ForgotPasswordView *views.View
	ResetPasswordView  *views.View
	us                 models.UserService
	emailer            *email.Client
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

type ResetPasswordForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

func NewUsers(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		NewView:            views.NewView("bootstrap", "users/new"),
		LoginView:          views.NewView("bootstrap", "users/login"),
		ForgotPasswordView: views.NewView("bootstrap", "users/forgot_password"),
		ResetPasswordView:  views.NewView("bootstrap", "users/reset_password"),
		us:                 us,
		emailer:            emailer,
	}
}

// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var signupForm SignupForm

	err := parseForm(r, &signupForm)
	if err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	user := models.User{
		Username: signupForm.Username,
		Email:    signupForm.Email,
		Password: signupForm.Password,
	}

	err = u.us.Create(&user)
	if err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}

	u.emailer.Welcome(user.Username, user.Email)
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
	var vd views.Data
	var loginForm LoginForm
	if err := parseForm(r, &loginForm); err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}

	user, err := u.us.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError(models.ErrNoUserWithEmail)
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}

	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
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

// POST /logout
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	// First expire the user's cookie
	cookie := http.Cookie{
		Name:     "authToken",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// Then we update the user with a new remember token
	user := middlewares.LookUpUserFromContext(r.Context())
	// We are ignoring errors for now because they are
	// unlikely, and even if they do occur we can't recover
	// now that the user doesn't have a valid cookie
	token, _ := utils.GenerateToken()
	user.Token = token
	u.us.Update(user)
	// Finally send the user to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// POST /forgot
func (u *Users) InitiateReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPasswordForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ForgotPasswordView.Render(w, r, vd)
		return
	}

	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPasswordView.Render(w, r, vd)
		return
	}

	// Email the user their password reset token.
	err = u.emailer.ResetPassword(form.Email, token)
	if err != nil {
		vd.SetAlert(err)
		u.ForgotPasswordView.Render(w, r, vd)
		return
	}

	views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Instructions for resetting your password have been emailed to you.",
	})
}

// GET /reset
func (u *Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPasswordForm
	vd.Yield = &form
	err := parseURLParams(r, &form)
	if err != nil {
		vd.SetAlert(err)
	}
	u.ResetPasswordView.Render(w, r, vd)
}

// POST /reset
func (u *Users) CompleteReset(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ResetPasswordForm
	vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		u.ResetPasswordView.Render(w, r, vd)
		return
	}

	user, err := u.us.CompleteReset(form.Token, form.Password)
	if err != nil {
		vd.SetAlert(err)
		u.ResetPasswordView.Render(w, r, vd)
		return
	}

	u.signIn(w, user)
	views.RedirectAlert(w, r, "/galleries", http.StatusFound, views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Your password has been reset and you have been logged in!",
	})
}
