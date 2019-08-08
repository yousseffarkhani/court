package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"contactMail"
	"courtdb"
	"handlers"
	"middlewares"
	"model"
	"views"

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
	server.store = store
	router := newRouter(server)
	server.Handler = router
	return server, nil
}

func newRouter(server *BasketServer) *mux.Router {
	router := mux.NewRouter()

	authorizationMiddleware := middlewares.Authorization
	loggedMiddleware := middlewares.Logged

	router.Handle("/", loggedMiddleware.ThenFunc(server.indexHandler)).Methods(http.MethodGet)

	router.Handle("/contact", loggedMiddleware.ThenFunc(handlers.ContactGetHandler)).Methods(http.MethodGet)
	router.HandleFunc("/contact", ContactPostHandler).Methods(http.MethodPost)

	router.Handle("/login", loggedMiddleware.ThenFunc(handlers.LoginHandler)).Methods(http.MethodGet)
	router.HandleFunc("/signin", server.signinHandler).Methods(http.MethodPost)
	router.Handle("/signup", loggedMiddleware.ThenFunc(handlers.SignupHandler)).Methods(http.MethodGet)
	router.HandleFunc("/signup", server.signupHandler).Methods(http.MethodPost)
	router.HandleFunc("/logout", handlers.Logout).Methods(http.MethodGet)

	router.Handle("/court/new", authorizationMiddleware.ThenFunc(handlers.NewCourtHandler)).Methods(http.MethodGet)
	router.Handle("/court/new", authorizationMiddleware.ThenFunc(server.newCourtHandler)).Methods(http.MethodPost)
	router.Handle("/court/{id}", loggedMiddleware.ThenFunc(server.getCourtHandler)).Methods(http.MethodGet)
	router.HandleFunc("/api/courts", server.getAllCourtsHandler).Methods(http.MethodGet)

	router.HandleFunc("/court/{id}/comment", server.getComments).Methods(http.MethodGet)
	// router.Handle("/court/{id}/comment", authorizationMiddleware.ThenFunc(server.addComment)).Methods(http.MethodPost)
	router.HandleFunc("/court/{id}/comment/new", server.addComment).Methods(http.MethodPost)
	// router.Handle("/court/{id}/comment", authorizationMiddleware.ThenFunc(server.addComment)).Methods(http.MethodPost) // TODO: add Authorization (modifier le bouton en onClick action)
	router.Handle("/court/{id}/comment/delete", authorizationMiddleware.ThenFunc(server.deleteComment)).Methods(http.MethodPost) // TODO: user ID verification
	router.Handle("/court/{id}/comment/update", authorizationMiddleware.ThenFunc(server.updateComment)).Methods(http.MethodPost) // TODO: user ID verification

	// router.Handle("/admin", authorizationMiddleware.ThenFunc(adminHandler)).Methods(http.MethodGet) // TODO: Créer une page admin pour pouvoir accepter les soumissions des nouveaux terrains

	// router.HandleFunc("/court/{id}", DeleteCourt).Methods(http.MethodDelete)
	// router.HandleFunc("/court/{id}", UpdateCourt).Methods(http.MethodPut)

	router.HandleFunc("/static/css/style.css", serveCss).Methods(http.MethodGet)
	router.HandleFunc("/static/js/main.js", serveJs).Methods(http.MethodGet)

	router.PathPrefix("/").HandlerFunc(redirectToIndex).Methods("GET")
	return router
}

func (server *BasketServer) addComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var comment model.Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	comment.CourtID = id
	if comment.Message == "" { // TODO: Vérification de l'input
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = server.store.AddComment(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (server *BasketServer) getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", JsonContentType)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	comments, err := server.store.GetComments(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&comments)
}

func (server *BasketServer) deleteComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = server.store.DeleteComment(comment.ID)
	if err != nil {
		if err.Error() == "Comment doesn't exist." {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (server *BasketServer) updateComment(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	// TODO: Vérifier les inputs
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = server.store.UpdateComment(comment)
	if err != nil {
		if err.Error() == "Comment doesn't exist." {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (server *BasketServer) signupHandler(w http.ResponseWriter, r *http.Request) {
	userInput, inputErrors := getInput(w, r)

	templateData := model.TemplateData{
		ActionDone: false,
		Errors:     inputErrors,
		Ressource:  userInput.Username,
	}
	if len(inputErrors) > 0 {
		err := views.Pages["signup"].Render(w, r, templateData)
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
		errors := make(map[string]string)
		errors["Error"] = err.Error()
		templateData.Errors = errors
		err := views.Pages["signup"].Render(w, r, templateData)
		if err != nil {
			http.Redirect(w, r, "/signup", http.StatusFound)
		}
		return
	}

	middlewares.SetJwtCookie(w, newUser.Username)
	templateData.ActionDone = true
	err = views.Pages["signup"].Render(w, r, templateData)
	if err != nil {
		http.Redirect(w, r, "/signup", http.StatusFound)
	}
}

func (server *BasketServer) signinHandler(w http.ResponseWriter, r *http.Request) {
	userInput, inputErrors := getInput(w, r)
	if len(inputErrors) > 0 {
		err := views.Pages["login"].Render(w, r, inputErrors)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}
	user, err := server.store.GetUser(userInput.Username)
	if err != nil {
		errors := [1]string{err.Error()}
		err := views.Pages["login"].Render(w, r, errors)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		errors := [1]string{"Wrong password"}
		err := views.Pages["login"].Render(w, r, errors)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}

	middlewares.SetJwtCookie(w, user.Username)
	redirectToIndex(w, r)
}

func (server *BasketServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContentType)
	courts := server.store.GetAllCourts()
	courtsBytes, err := json.MarshalIndent(courts, "", " ")
	CheckError(err)
	err = views.Pages["index"].Render(w, r, string(courtsBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* Court Handlers */
func (server *BasketServer) getAllCourtsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", JsonContentType)
	courts := server.store.GetAllCourts()
	CheckError(json.NewEncoder(w).Encode(courts))
}

func (server *BasketServer) getCourtHandler(w http.ResponseWriter, r *http.Request) {
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
	err = views.Pages["courtDetails"].Render(w, r, court)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *BasketServer) newCourtHandler(w http.ResponseWriter, r *http.Request) {
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

	templateData := model.TemplateData{
		ActionDone: false,
		Errors:     errors,
		Ressource:  courtInput,
	}

	if len(errors) > 0 {
		err := views.Pages["newCourt"].Render(w, r, templateData)
		if err != nil {
			http.Redirect(w, r, "/signup", http.StatusFound)
		}
		return
	}

	err := server.store.AddCourt(courtInput)
	if err != nil {
		errors := make(map[string]string)
		errors["Error"] = err.Error()
		templateData.Errors = errors
		err := views.Pages["newCourt"].Render(w, r, templateData)
		if err != nil {
			http.Redirect(w, r, "/courts/new", http.StatusFound)
		}
		return
	}

	templateData.ActionDone = true
	err = views.Pages["newCourt"].Render(w, r, templateData)
	if err != nil {
		http.Redirect(w, r, "/courts/new", http.StatusFound)
	}
}

/* Contact handler */
func ContactPostHandler(w http.ResponseWriter, r *http.Request) {
	newContact := contactMail.Contact{
		Name:    strings.TrimSpace(r.FormValue("name")),
		Subject: strings.TrimSpace(r.FormValue("subject")),
		Email:   strings.TrimSpace(r.FormValue("email")),
		Message: strings.TrimSpace(r.FormValue("message")),
	}

	errors := extractEmptyFieldErrors(newContact)

	templateData := model.TemplateData{
		ActionDone: false,
		Errors:     errors,
		Ressource:  newContact,
	}

	if len(errors) > 0 {
		err := views.Pages["contact"].Render(w, r, templateData)
		if err != nil {
			http.Redirect(w, r, "/contact", http.StatusFound)
		}
		return
	}

	err := contactMail.SendMail(newContact)
	if err != nil {
		errors := map[string]string{
			"errMail": err.Error(),
		}
		templateData.Errors = errors
		err := views.Pages["contact"].Render(w, r, templateData)
		if err != nil {
			http.Redirect(w, r, "/contact", http.StatusFound)
		}
		return
	}

	fmt.Println("Email sent")
	templateData.ActionDone = true
	err = views.Pages["contact"].Render(w, r, templateData)
	if err != nil {
		http.Redirect(w, r, "/contact", http.StatusFound)
	}
}

/* Utils */
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
