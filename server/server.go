package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yousseffarkhani/court/courtdb"
	"github.com/yousseffarkhani/court/views"
)

const JsonContentType = "application/json"
const HtmlContentType = "text/html"

var index *views.View

// var contact *views.View
var courtDetails *views.View

type BasketServer struct {
	store *courtdb.CourtStore
	http.Handler
	// courtsCache []courtdb.Court
}

func NewBasketServer(store *courtdb.CourtStore) (*BasketServer, error) {
	server := new(BasketServer)

	initializeViews()
	server.store = store
	router := newRouter(server)
	server.Handler = router
	// server.courtsCache = server.store.GetAllCourts() // TODO: Refresh cache every week
	return server, nil
}

func initializeViews() {
	index = views.NewView("main", "views/index.html")
	courtDetails = views.NewView("main", "views/court.html")
}

func newRouter(server *BasketServer) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", server.indexHandler).Methods(http.MethodGet)
	router.HandleFunc("/courts/{id}", server.courtHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/courts", server.getAllCourts).Methods(http.MethodGet)
	// router.HandleFunc("/courts/{id}", DeleteCourt).Methods(http.MethodDelete)
	// router.HandleFunc("/courts/{id}", UpdateCourt).Methods(http.MethodPut)

	router.HandleFunc("/static/css/style.css", serveCss).Methods(http.MethodGet)
	router.HandleFunc("/static/js/main.js", serveJs).Methods(http.MethodGet)
	router.PathPrefix("/").HandlerFunc(redirect).Methods("GET")
	return router
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/css/style.css")
}

func serveJs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/js/main.js")
}

func (server BasketServer) getAllCourts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", JsonContentType)
	courts := server.store.GetAllCourts()
	CheckError(json.NewEncoder(w).Encode(courts))
}

func (server BasketServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContentType)
	courts := server.store.GetAllCourts()
	courtsBytes, err := json.MarshalIndent(courts, "", " ")
	CheckError(err)
	err = index.Render(w, string(courtsBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server BasketServer) courtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContentType)
	vars := mux.Vars(r)
	id := vars["id"]
	court := server.store.GetCourt(id)
	if court.ID == 0 {
		redirect(w, r)
	}
	err := courtDetails.Render(w, court)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
