package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yousseffarkhani/court/model"
	"github.com/yousseffarkhani/court/views"
)

/* Court Handlers */
func (server *BasketServer) APIGetAllCourts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", JsonContentType)
	courts := server.database.GetAllCourts()
	if err := json.NewEncoder(w).Encode(courts); err != nil {
		log.Fatalln(err)
	}
}

func (server *BasketServer) APIGetCourt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContentType)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		redirectToIndex(w, r)
	}
	court := server.database.GetCourt(id)
	if court.ID == 0 {
		redirectToIndex(w, r)
	}
	err = views.Pages["courtDetails"].Render(w, r, court)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *BasketServer) APINewCourt(w http.ResponseWriter, r *http.Request) {
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

	err := server.database.AddCourt(courtInput)
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
