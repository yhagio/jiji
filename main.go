package main

import (
	"flag"
	"fmt"
	"jiji/controllers"
	"jiji/middlewares"
	"jiji/models"
	"jiji/utils"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
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
	)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

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
	usersCtrl := controllers.NewUsers(services.User)
	galleriesCtrl := controllers.NewGalleries(services.Gallery, services.Image, r)

	// ********* Middlewares *********
	userMW := middlewares.User{
		UserService: services.User,
	}
	requireUserMW := middlewares.RequireUser{}

	// CSRF
	bytes, err := utils.GenerateRandomBytes(32)
	if err != nil {
		panic(err)
	}
	csrfMW := csrf.Protect(bytes, csrf.Secure(config.IsProd()))

	// ********* Static page *********
	r.Handle("/", staticCtrl.HomeView).Methods("GET")
	r.Handle("/contact", staticCtrl.ContactView).Methods("GET")

	// ********* Users *********
	r.HandleFunc("/signup", usersCtrl.New).Methods("GET")
	r.HandleFunc("/signup", usersCtrl.Create).Methods("POST")
	r.Handle("/login", usersCtrl.LoginView).Methods("GET")
	r.HandleFunc("/login", usersCtrl.Login).Methods("POST")

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
