package court

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

var tpl *template.Template

func init() {
	fp := path.Join("assets/templates", "court.html")
	tpl = template.Must(template.ParseFiles(fp))
}

type Courts []Court

type Court struct {
	Name           string `json:"nom"`
	Url            string `json:"url"`
	Adress         string `json:"adresse"`
	Arrondissement string `json:"arrondissement"`
	Longitude      string `json:"longitude"`
	Lattitude      string `json:"lattitude"`
	Dimensions     string `json:"dimensions"`
	Revetement     string `json:"revetement"`
	Decouvert      string `json:"decouvert"`
	Eclairage      string `json:"eclairage"`
}

func NewHandler() http.Handler {
	h := handler{
		court:  Court{},
		t:      tpl,
		pathFc: defaultPathFc}
	return h
}

type handler struct {
	court  Court
	t      *template.Template
	pathFc func(r *http.Request) string
}

func defaultPathFc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	fmt.Println(path)
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = h.pathFc(r)
	/* [{"nom":"TEP NEUVE SAINT PIERRE","url":"http://www.webvilles.net/sports/activites/140062/terrain-de-mini--basket-ball-paris.php","adresse":"Tep Neuve Saint Pierre Paris 5 Rue Neuve St Pierre 75004 Paris","arrondissement":"75004","longitude":"2.36295000","lattitude":"48.85351000","dimensions":"160.00","revetement":"Synthétique (hors gazon)","decouvert":"Découvert","eclairage":"oui"} */

	mockCourt := Court{
		Adress:         "Tep Neuve Saint Pierre Paris 5 Rue Neuve St Pierre 75004 Paris",
		Arrondissement: "75004",
		Decouvert:      "Découvert",
		Dimensions:     "160.00",
		Eclairage:      "oui",
		Lattitude:      "48.85351000",
		Longitude:      "2.36295000",
		Name:           "TEP NEUVE SAINT PIERRE",
		Revetement:     "Synthétique (hors gazon)",
		Url:            "http://www.webvilles.net/sports/activites/140062/terrain-de-mini--basket-ball-paris.php",
	} // TODO: Il faut query la BDD avec le path
	h.court = mockCourt
	err := h.t.Execute(w, mockCourt)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Something went wrong ...", http.StatusInternalServerError)
	}
	return
}

func GetCourt() {

}
