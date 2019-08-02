package server

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"courtdb"
	"middlewares"
	"model"
	"views"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

const JsonContentType = "application/json"
const HtmlContentType = "text/html"

type BasketServer struct {
	store *courtdb.CourtStore
	http.Handler
}

func NewBasketServer(store *courtdb.CourtStore) (*BasketServer, error) {
	server := new(BasketServer)

	initializeViews()
	server.store = store
	router := newRouter(server)
	server.Handler = router
	return server, nil
}

// var contact *views.View
var index *views.View
var courtDetails *views.View
var login *views.View
var signup *views.View
var newCourt *views.View

func initializeViews() {
	index = views.NewView("main", "views/index.html")
	login = views.NewView("main", "views/login.html")
	signup = views.NewView("main", "views/signup.html")
	courtDetails = views.NewView("main", "views/court.html")
	newCourt = views.NewView("main", "views/newCourt.html")
}

func newRouter(server *BasketServer) *mux.Router {
	router := mux.NewRouter()

	authorizationMiddleware := middlewares.Authorization
	loggedMiddleware := middlewares.Logged

	router.Handle("/", loggedMiddleware.ThenFunc(server.indexHandler)).Methods(http.MethodGet)

	router.Handle("/login", loggedMiddleware.ThenFunc(loginHandler)).Methods(http.MethodGet)
	router.HandleFunc("/signin", server.signinHandler).Methods(http.MethodPost)
	router.Handle("/signup", loggedMiddleware.ThenFunc(signupHandler)).Methods(http.MethodGet)
	router.HandleFunc("/signup", server.signupHandler).Methods(http.MethodPost)
	router.HandleFunc("/logout", Logout).Methods(http.MethodGet)

	router.Handle("/courts/new", authorizationMiddleware.ThenFunc(newCourtHandler)).Methods(http.MethodGet)
	router.Handle("/courts/new", authorizationMiddleware.ThenFunc(server.newCourtHandler)).Methods(http.MethodPost)
	router.Handle("/courts/{id}", loggedMiddleware.ThenFunc(server.getCourtHandler)).Methods(http.MethodGet)
	router.HandleFunc("/api/courts", server.getAllCourtsHandler).Methods(http.MethodGet)

	// router.Handle("/admin", authorizationMiddleware.ThenFunc(adminHandler)).Methods(http.MethodGet) // Créer une page admin pour pouvoir accepter les soumissions des nouveaux terrains

	// router.HandleFunc("/courts/{id}", DeleteCourt).Methods(http.MethodDelete)
	// router.HandleFunc("/courts/{id}", UpdateCourt).Methods(http.MethodPut)

	router.HandleFunc("/static/css/style.css", serveCss).Methods(http.MethodGet)
	router.HandleFunc("/static/js/main.js", serveJs).Methods(http.MethodGet)

	router.PathPrefix("/").HandlerFunc(redirectToIndex).Methods("GET")
	return router
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	userLogged, _ := context.Get(r, "userLogged").(bool)
	if userLogged {
		redirectToIndex(w, r)
	} else {
		err := login.Render(w, r, nil)
		if err != nil {

			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	expired := time.Now().Add(time.Minute * -5)
	http.SetCookie(w, &http.Cookie{
		Name:    "Token",
		Value:   "",
		Expires: expired,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	userLogged, _ := context.Get(r, "userLogged").(bool)
	if userLogged {
		redirectToIndex(w, r)
	} else {
		err := signup.Render(w, r, nil)
		if err != nil {

			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func (server BasketServer) signupHandler(w http.ResponseWriter, r *http.Request) {
	userInput, inputErrors := getInput(w, r)
	if len(inputErrors) > 0 {
		err := signup.Render(w, r, inputErrors)
		if err != nil {
			http.Redirect(w, r, "/signup", http.StatusFound)
		}
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 8)
	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	newUser := model.User{
		Username: userInput.Username,
		Password: string(hashedPassword),
	}

	err = server.store.AddUser(newUser)
	if err != nil {
		errors := [1]string{err.Error()}
		err := signup.Render(w, r, errors)
		if err != nil {

			http.Redirect(w, r, "/signup", http.StatusFound)
		}
		return
	}
	middlewares.SetJwtCookie(w, newUser.Username)
	redirectToIndex(w, r)
}

func (server BasketServer) signinHandler(w http.ResponseWriter, r *http.Request) {
	userInput, inputErrors := getInput(w, r)
	if len(inputErrors) > 0 {
		err := login.Render(w, r, inputErrors)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}
	user, err := server.store.GetUser(userInput.Username)
	if err != nil {
		errors := [1]string{err.Error()}
		err := login.Render(w, r, errors)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		errors := [1]string{"Wrong password"}
		err := login.Render(w, r, errors)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}

	middlewares.SetJwtCookie(w, user.Username)
	redirectToIndex(w, r)
}

func (server BasketServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContentType)
	courts := server.store.GetAllCourts()
	courtsBytes, err := json.MarshalIndent(courts, "", " ")
	CheckError(err)
	err = index.Render(w, r, string(courtsBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* Court Handlers */
func (server BasketServer) getAllCourtsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", JsonContentType)
	courts := server.store.GetAllCourts()
	CheckError(json.NewEncoder(w).Encode(courts))
}

func (server BasketServer) getCourtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContentType)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		redirectToIndex(w, r)
	}
	court := server.store.GetCourt(id)
	if court.ID == 0 {
		redirectToIndex(w, r)
	}
	err = courtDetails.Render(w, r, court)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newCourtHandler(w http.ResponseWriter, r *http.Request) {
	err := newCourt.Render(w, r, nil)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (server BasketServer) newCourtHandler(w http.ResponseWriter, r *http.Request) {
	courtInput := model.Court{
		Name:           strings.TrimSpace(r.FormValue("name")),
		Url:            "court.Url",
		Adress:         strings.TrimSpace(r.FormValue("adress")),
		Arrondissement: "court.Arrondissement",
		Longitude:      "court.Longitude",
		Lattitude:      "court.Lattitude",
		Dimensions:     "court.Dimensions",
		Revetement:     "court.Revetement",
		Decouvert:      "court.Decouvert",
		Eclairage:      "court.Eclairage",
	}
	errors := extractEmptyFieldErrors(courtInput)

	if len(errors) > 0 {
		err := newCourt.Render(w, r, errors)
		if err != nil {
			http.Redirect(w, r, "/signup", http.StatusFound)
		}
		return
	}

	err := server.store.AddCourt(courtInput)
	if err != nil {
		errors := [1]string{err.Error()}
		err := newCourt.Render(w, r, errors)
		if err != nil {
			http.Redirect(w, r, "/courts/new", http.StatusFound)
		}
		return
	}
	redirectToIndex(w, r) // TODO: Display success message and stay on same page
}

func extractEmptyFieldErrors(data interface{}) map[string]string {
	errors := make(map[string]string)
	value := reflect.ValueOf(data)
	typeOfData := value.Type()
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			if value.Field(i).String() == "" {
				errors[typeOfData.Field(i).Name] = "Empty field"
			}
		}
	}
	return errors
}

/* Utils */
func getInput(w http.ResponseWriter, r *http.Request) (model.User, map[string]string) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))

	userInput := model.User{
		Username: username,
		Password: password,
	}

	errors := extractEmptyFieldErrors(userInput)

	if len(errors) > 0 {
		return userInput, errors
	}

	return userInput, nil
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/css/style.css")
}

func serveJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/js/main.js")
}
