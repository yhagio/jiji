package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jiji/middlewares"
	"jiji/models"
	"net/http"
	"time"

	"golang.org/x/oauth2"

	"github.com/gorilla/csrf"
)

type OAuth struct {
	os           models.OAuthService
	dropboxOAuth *oauth2.Config
}

func NewOAuth(os models.OAuthService, dropboxConfig *oauth2.Config) *OAuth {
	return &OAuth{
		os:           os,
		dropboxOAuth: dropboxConfig,
	}
}

func (oa *OAuth) DropboxConnect(w http.ResponseWriter, r *http.Request) {
	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	url := oa.dropboxOAuth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (oa *OAuth) DropboxCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state is provided", http.StatusBadRequest)
		return
	}

	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)

	code := r.FormValue("code")
	token, err := oa.dropboxOAuth.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := middlewares.LookUpUserFromContext(r.Context())
	exist, err := oa.os.Find(user.ID, models.OAuthDropbox)
	if err == models.ErrNotFound {
		// Nothing to do
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// Delete the existing one
		oa.os.Delete(exist.ID)
	}

	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: models.OAuthDropbox,
	}
	err = oa.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v", token)

	fmt.Fprintln(w, "code: ", r.FormValue("code"), " state: ", r.FormValue("state"))
}

func (oa *OAuth) DropboxTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := r.FormValue("path")

	user := middlewares.LookUpUserFromContext(r.Context())
	userOAuth, err := oa.os.Find(user.ID, models.OAuthDropbox)
	if err != nil {
		panic(err)
	}
	token := userOAuth.Token

	data := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	client := oa.dropboxOAuth.Client(context.TODO(), &token)
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.dropboxapi.com/2/files/list_folder",
		bytes.NewReader(dataBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	io.Copy(w, res.Body)
}
