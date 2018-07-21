package main

import (
	"flag"
	"fmt"
	"jiji/controllers"
	"jiji/email"
	"jiji/middlewares"
	"jiji/models"
	"jiji/utils"
	"net/http"

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
	oauthsCtrl := controllers.NewOAuth(services.OAuth, dropboxOAuth)

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

	// ********* OAuth *********
	r.HandleFunc("/oauth/dropbox/connect", requireUserMW.ApplyFunc(oauthsCtrl.DropboxConnect))
	r.HandleFunc("/oauth/dropbox/callback", requireUserMW.ApplyFunc(oauthsCtrl.DropboxCallback))
	r.HandleFunc("/oauth/dropbox/test", requireUserMW.ApplyFunc(oauthsCtrl.DropboxTest))

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
