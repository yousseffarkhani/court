package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/yousseffarkhani/court/database"
	"github.com/yousseffarkhani/court/middlewares"
	"github.com/yousseffarkhani/court/model"
	"github.com/yousseffarkhani/court/views"

	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

const JsonContentType = "application/json"
const HtmlContentType = "text/html"

type BasketServer struct {
	database *database.CourtStore
	http.Handler
}

func NewBasketServer(database *database.CourtStore) (*BasketServer, error) {
	server := new(BasketServer)
	server.database = database
	router := newRouter(server)
	server.Handler = router
	return server, nil
}

func newRouter(server *BasketServer) *mux.Router {
	router := mux.NewRouter()

	authorizationMiddleware := middlewares.Authorization
	loggedMiddleware := middlewares.Logged

	router.Handle("/", loggedMiddleware.ThenFunc(server.indexHandler)).Methods(http.MethodGet)

	router.Handle("/contact", loggedMiddleware.ThenFunc(getContactHandler)).Methods(http.MethodGet)
	router.HandleFunc("/contact", postContactHandler).Methods(http.MethodPost)

	router.Handle("/login", loggedMiddleware.ThenFunc(LoginHandler)).Methods(http.MethodGet)
	router.Handle("/signup", loggedMiddleware.ThenFunc(getSignupHandler)).Methods(http.MethodGet)
	router.HandleFunc("/signup", server.postSignupHandler).Methods(http.MethodPost)
	router.HandleFunc("/signin", server.signinHandler).Methods(http.MethodPost)
	router.HandleFunc("/logout", Logout).Methods(http.MethodGet)

	router.Handle("/court/new", authorizationMiddleware.ThenFunc(NewCourtHandler)).Methods(http.MethodGet)
	router.Handle("/court/new", authorizationMiddleware.ThenFunc(server.APINewCourt)).Methods(http.MethodPost)
	router.Handle("/court/{id}", loggedMiddleware.ThenFunc(server.APIGetCourt)).Methods(http.MethodGet)
	router.HandleFunc("/api/courts", server.APIGetAllCourts).Methods(http.MethodGet)

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
	message := strings.TrimSpace(r.FormValue("message"))

	comment := model.Comment{
		Message: message,
	}

	errors := extractEmptyFieldErrors(comment)

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment.CourtID = id
	err = server.database.AddComment(comment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/court/%d", id), http.StatusFound)
	// w.WriteHeader(http.StatusCreated)
}

func (server *BasketServer) getComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", JsonContentType)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	comments, err := server.database.GetComments(id)
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
	err = server.database.DeleteComment(comment.ID)
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
	err = server.database.UpdateComment(comment)
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

func (server *BasketServer) signinHandler(w http.ResponseWriter, r *http.Request) {
	userInput, inputErrors := getInput(w, r)
	if len(inputErrors) > 0 {
		err := views.Pages["login"].Render(w, r, inputErrors)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}
	user, err := server.database.GetUser(userInput.Username)
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

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/css/style.css")
}

func serveJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/js/main.js")
}
