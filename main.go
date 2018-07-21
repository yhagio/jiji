package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"jiji/controllers"
	"jiji/email"
	"jiji/middlewares"
	"jiji/models"
	"jiji/utils"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

func main() {
	// ********* Create User Service *********
	configRequired := flag.Bool("prod", false, "Provide this flag in production. "+
		"This ensures that a .config file is provided before the application starts.")
	flag.Parse()

	config := LoadConfig(*configRequired)
	dbConfig := config.Database
	services, err := models.NewServices(
		models.WithGorm(dbConfig.Dialect(), dbConfig.ConnectionInfo()),
		models.WithLogMode(!config.IsProd()),
		models.WithUser(config.Pepper, config.HMACKey),
		models.WithGallery(),

		models.WithImage(),
		models.WithOAuth(),
	)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	services.AutoMigrate()
	// services.DestructiveReset()

	// Mailgun config
	mailgunConfig := config.Mailgun
	emailer := email.NewClient(
		email.WithSender("JIJI Support", "support@"+mailgunConfig.Domain),
		email.WithMailgun(mailgunConfig.Domain, mailgunConfig.APIKey, mailgunConfig.PublicAPIKey),
	)

	r := mux.NewRouter()

	// ********* Assets *********
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// ********* Image handler *********
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// ********* Defines controllers *********
	staticCtrl := controllers.NewStatic()
	usersCtrl := controllers.NewUsers(services.User, emailer)
	galleriesCtrl := controllers.NewGalleries(services.Gallery, services.Image, r)

	// ********* Middlewares *********
	userMW := middlewares.User{
		UserService: services.User,
	}
	requireUserMW := middlewares.RequireUser{}

	// CSRF
	generatedBytes, err := utils.GenerateRandomBytes(32)
	if err != nil {
		panic(err)
	}
	// If get CSRF token is invalid error, app is ruuning on localhost or non-https
	csrfMW := csrf.Protect(generatedBytes, csrf.Secure(config.IsProd()))

	// OAuth dropbox
	dropboxOAuth := &oauth2.Config{
		ClientID:     config.Dropbox.ID,
		ClientSecret: config.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Dropbox.AuthURL,
			TokenURL: config.Dropbox.TokenURL,
		},
		RedirectURL: "http://localhost:3000/oauth/dropbox/callback",
	}

	dropboxRedirect := func(w http.ResponseWriter, r *http.Request) {
		state := csrf.Token(r)
		cookie := http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		url := dropboxOAuth.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
	}

	dropboxCallback := func(w http.ResponseWriter, r *http.Request) {

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
		token, err := dropboxOAuth.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user := middlewares.LookUpUserFromContext(r.Context())
		exist, err := services.OAuth.Find(user.ID, models.OAuthDropbox)
		if err == models.ErrNotFound {
			// Nothing to do
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			// Delete the existing one
			services.OAuth.Delete(exist.ID)
		}

		userOAuth := models.OAuth{
			UserID:  user.ID,
			Token:   *token,
			Service: models.OAuthDropbox,
		}
		err = services.OAuth.Create(&userOAuth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%+v", token)

		fmt.Fprintln(w, "code: ", r.FormValue("code"), " state: ", r.FormValue("state"))
	}

	dropboxQuery := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		path := r.FormValue("path")

		user := middlewares.LookUpUserFromContext(r.Context())
		userOAuth, err := services.OAuth.Find(user.ID, models.OAuthDropbox)
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

		client := dropboxOAuth.Client(context.TODO(), &token)
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

	r.HandleFunc("/oauth/dropbox/connect", requireUserMW.ApplyFunc(dropboxRedirect))
	r.HandleFunc("/oauth/dropbox/callback", requireUserMW.ApplyFunc(dropboxCallback))
	r.HandleFunc("/oauth/dropbox/test", requireUserMW.ApplyFunc(dropboxQuery))

	// ********* Static page *********
	r.Handle("/", staticCtrl.HomeView).Methods("GET")
	r.Handle("/contact", staticCtrl.ContactView).Methods("GET")

	// ********* Users *********
	r.HandleFunc("/signup", usersCtrl.New).Methods("GET")
	r.HandleFunc("/signup", usersCtrl.Create).Methods("POST")
	r.Handle("/login", usersCtrl.LoginView).Methods("GET")
	r.HandleFunc("/login", usersCtrl.Login).Methods("POST")
	r.Handle("/logout", requireUserMW.ApplyFunc(usersCtrl.Logout)).Methods("POST")

	r.Handle("/forgot", usersCtrl.ForgotPasswordView).Methods("GET")
	r.HandleFunc("/forgot", usersCtrl.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", usersCtrl.ResetPassword).Methods("GET")
	r.HandleFunc("/reset", usersCtrl.CompleteReset).Methods("POST")

	// ********* Galleries *********
	r.Handle("/galleries", requireUserMW.ApplyFunc(galleriesCtrl.GetAllByUser)).Methods("GET")
	r.Handle("/galleries/new", requireUserMW.Apply(galleriesCtrl.New)).Methods("GET")
	r.Handle("/galleries", requireUserMW.ApplyFunc(galleriesCtrl.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesCtrl.Show).Methods("GET").Name(controllers.ShowGallery)

	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMW.ApplyFunc(galleriesCtrl.Edit)).Methods("GET").Name(controllers.EditGallery)

	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMW.ApplyFunc(galleriesCtrl.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMW.ApplyFunc(galleriesCtrl.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMW.ApplyFunc(galleriesCtrl.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMW.ApplyFunc(galleriesCtrl.DeleteImage)).Methods("POST")

	fmt.Printf("Starting the server on localhost:%d...\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), csrfMW(userMW.Apply(r)))
}
