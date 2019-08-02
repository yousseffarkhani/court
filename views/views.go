package views

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/context"
)

const LayoutDir = "views/layouts"

type View struct {
	Template *template.Template
	Layout   string
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
