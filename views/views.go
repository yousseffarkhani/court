package views

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/context"
	"github.com/yousseffarkhani/court/model"
)

const LayoutDir = "views/layouts"

type View struct {
	Template *template.Template
	Layout   string
}

var Pages = make(map[string]*View)

func init() {
	Pages["index"] = NewView("main", "views/index.html")
	Pages["login"] = NewView("main", "views/login.html")
	Pages["signup"] = NewView("main", "views/signup.html")
	Pages["courtDetails"] = NewView("main", "views/court.html")
	Pages["newCourt"] = NewView("main", "views/newCourt.html")
	Pages["contact"] = NewView("main", "views/contact.html")
}

func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	tmpl := template.Must(template.ParseFiles(files...))
	return &View{
		Template: tmpl,
		Layout:   layout,
	}
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "/*.html")
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	userLogged, _ := context.Get(r, "userLogged").(bool)
	renderingData := RenderingData{
		UserLogged: userLogged,
		Data:       data,
	}
	err := v.Template.ExecuteTemplate(w, v.Layout, renderingData)
	return err
}

type RenderingData struct {
	Data       interface{}
	UserLogged bool
}

func RenderIndex(w http.ResponseWriter, r *http.Request, courts model.Courts) error {
	courtsBytes, err := json.MarshalIndent(courts, "", " ")
	if err != nil {
		return fmt.Errorf("Problem encoding to JSON, %v", err)
	}
	err = Pages["index"].Render(w, r, string(courtsBytes))
	if err != nil {
		return fmt.Errorf("Problem encoding to JSON, %v", err)
	}
	return nil
}
