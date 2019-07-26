package views

import (
	"net/http"
	"path/filepath"
	"text/template"
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

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "/*.html")
	if err != nil {
		panic(err)
	}
	return files
}
