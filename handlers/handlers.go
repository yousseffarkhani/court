package handlers

import (
	"contactMail"
	"model"
	"net/http"
	"time"
	"views"

	"github.com/gorilla/context"
)

/* Login/Logout handlers */
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	userLogged, _ := context.Get(r, "userLogged").(bool)
	if userLogged {
		redirectToIndex(w, r)
	} else {
		err := views.Pages["login"].Render(w, r, nil)
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

/* Signup handlers*/
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	userLogged, _ := context.Get(r, "userLogged").(bool)
	if userLogged {
		redirectToIndex(w, r)
	} else {
		err := views.Pages["signup"].Render(w, r, model.TemplateData{
			Ressource: "",
		})
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

/* Court handlers */
func NewCourtHandler(w http.ResponseWriter, r *http.Request) {
	err := views.Pages["newCourt"].Render(w, r, model.TemplateData{
		Ressource: model.Court{},
	})
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

/* Contact handlers */

func ContactGetHandler(w http.ResponseWriter, r *http.Request) {
	err := views.Pages["contact"].Render(w, r, model.TemplateData{
		Ressource: contactMail.Contact{},
	})
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

/* Util */
func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
