package controllers

import (
	"context"
	"fmt"
	"jiji/dbx"
	"jiji/middlewares"
	"jiji/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"golang.org/x/oauth2"

	"github.com/gorilla/csrf"
)

type OAuth struct {
	os      models.OAuthService
	configs map[string]*oauth2.Config
}

func NewOAuth(os models.OAuthService, configs map[string]*oauth2.Config) *OAuth {
	return &OAuth{
		os:      os,
		configs: configs,
	}
}

func (oa *OAuth) Connect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := oa.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (oa *OAuth) Callback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	oauthConfig, ok := oa.configs[service]
	if !ok {
		http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
		return
	}

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
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := middlewares.LookUpUserFromContext(r.Context())
	exist, err := oa.os.Find(user.ID, service)
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
		Service: service,
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
	vars := mux.Vars(r)
	service := vars["service"]
	// oauthConfig, ok := oa.configs[service]
	// if !ok {
	// 	http.Error(w, "Invalid OAuth2 Service", http.StatusBadRequest)
	// 	return
	// }

	r.ParseForm()
	path := r.FormValue("path")

	user := middlewares.LookUpUserFromContext(r.Context())
	userOAuth, err := oa.os.Find(user.ID, service)
	if err != nil {
		panic(err)
	}
	token := userOAuth.Token

	folders, files, err := dbx.GetList(token.AccessToken, path)
	if err != nil {
		panic(err)
	}

	// config := dropbox.Config{
	// 	Token: token.AccessToken,
	// }
	// dbx := files.New(config)

	// result, err := dbx.ListFolder(&files.ListFolderArg{
	// 	Path: path,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// for _, entry := range result.Entries {
	// 	switch meta := entry.(type) {
	// 	case *files.FolderMetadata:
	// 		fmt.Fprintln(w, "FolderMetadata: ", meta)
	// 	case *files.FileMetadata:
	// 		fmt.Fprintln(w, "FileMetadata: ", meta)
	// 	}
	// }

	// data := struct {
	// 	Path string `json:"path"`
	// }{
	// 	Path: path,
	// }
	// dataBytes, err := json.Marshal(data)
	// if err != nil {
	// 	panic(err)
	// }

	// client := oauthConfig.Client(context.TODO(), &token)
	// req, err := http.NewRequest(
	// 	http.MethodPost,
	// 	"https://api.dropboxapi.com/2/files/list_folder",
	// 	bytes.NewReader(dataBytes))
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Add("Content-Type", "application/json")
	// res, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer res.Body.Close()
	// io.Copy(w, res.Body)
}
