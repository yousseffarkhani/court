package server

import (
	"net/http"

	"github.com/yousseffarkhani/court/database"
	"github.com/yousseffarkhani/court/middlewares"

	"github.com/gorilla/mux"
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

	router.Handle("/", middlewares.Logged.ThenFunc(server.indexHandler)).Methods(http.MethodGet)

	router.Handle("/contact", middlewares.Logged.ThenFunc(getContactHandler)).Methods(http.MethodGet)
	router.HandleFunc("/contact", postContactHandler).Methods(http.MethodPost)

	router.Handle("/login", middlewares.Logged.ThenFunc(LoginHandler)).Methods(http.MethodGet)
	router.Handle("/signup", middlewares.Logged.ThenFunc(getSignupHandler)).Methods(http.MethodGet)
	router.HandleFunc("/signup", server.postSignupHandler).Methods(http.MethodPost)
	router.HandleFunc("/signin", server.signinHandler).Methods(http.MethodPost)
	router.HandleFunc("/logout", logout).Methods(http.MethodGet)

	router.Handle("/court/new", middlewares.Authorization.ThenFunc(getNewCourtHandler)).Methods(http.MethodGet)
	router.Handle("/court/new", middlewares.Authorization.ThenFunc(server.postNewCourtHandler)).Methods(http.MethodPost)
	router.Handle("/court/{id}", middlewares.Logged.ThenFunc(server.getCourtHandler)).Methods(http.MethodGet)

	/* API */
	router.HandleFunc("/api/courts", server.APIGetAllCourts).Methods(http.MethodGet)

	router.Handle("/court/{id}/comment", middlewares.Logged.ThenFunc(server.getComments)).Methods(http.MethodGet)
	router.Handle("/court/{id}/comment/new", middlewares.Authorization.ThenFunc(server.addComment)).Methods(http.MethodPost)
	router.Handle("/court/{id}/comment/{commentID}/delete", middlewares.Authorization.ThenFunc(server.deleteComment)).Methods(http.MethodPost)
	router.Handle("/court/{id}/comment/{commentID}/update", middlewares.Authorization.ThenFunc(server.updateComment)).Methods(http.MethodPost)

	// router.Handle("/admin", middlewares.Authorization.ThenFunc(adminHandler)).Methods(http.MethodGet) // TODO: Créer une page admin pour pouvoir accepter les soumissions des nouveaux terrains

	router.HandleFunc("/static/css/style.css", serveCss).Methods(http.MethodGet)
	router.HandleFunc("/static/js/main.js", serveJs).Methods(http.MethodGet)

	router.PathPrefix("/").HandlerFunc(redirectToIndex).Methods("GET")

	return router
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/css/style.css")
}

func serveJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/js/main.js")
}
