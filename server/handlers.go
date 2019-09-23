package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/yousseffarkhani/court/contactMail"
	"github.com/yousseffarkhani/court/middlewares"
	"github.com/yousseffarkhani/court/model"
	"github.com/yousseffarkhani/court/views"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func (server *BasketServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HtmlContentType)
	courts := server.database.GetAllCourts()
	err := views.RenderIndex(w, r, courts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* Contact handlers */
func getContactHandler(w http.ResponseWriter, r *http.Request) {
	err := views.Pages["contact"].Render(w, r, model.TemplateData{
		Ressource: contactMail.Contact{},
	})
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// TODO : Refactor because similar to postSignupHandler
func postContactHandler(w http.ResponseWriter, r *http.Request) {
	newContact := contactMail.Contact{
		Name:    extractFieldFromForm("name", r),
		Subject: extractFieldFromForm("subject", r),
		Email:   extractFieldFromForm("email", r),
		Message: extractFieldFromForm("message", r),
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

/* Signup handlers*/
func getSignupHandler(w http.ResponseWriter, r *http.Request) {
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

func (server *BasketServer) postSignupHandler(w http.ResponseWriter, r *http.Request) {
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

	err = server.database.AddUser(newUser)
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
	user, err := server.database.GetUser(userInput.Username)
	if err != nil {
		errors := map[string]string{"Error": err.Error()}
		err := views.Pages["login"].Render(w, r, errors)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		// TODO : Change error message because it gives too much info
		errors := map[string]string{"Error": "Wrong password"}
		err := views.Pages["login"].Render(w, r, errors)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		return
	}

	middlewares.SetJwtCookie(w, user.Username)
	redirectToIndex(w, r)
}

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

func logout(w http.ResponseWriter, r *http.Request) {
	expired := time.Now().Add(time.Minute * -5)
	http.SetCookie(w, &http.Cookie{
		Name:    "Token",
		Value:   "",
		Expires: expired,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

/* Court handlers */
func getNewCourtHandler(w http.ResponseWriter, r *http.Request) {
	err := views.Pages["newCourt"].Render(w, r, model.TemplateData{
		Ressource: model.Court{},
	})
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// TODO : Refactor because similar to postSignupHandler
func (server *BasketServer) postNewCourtHandler(w http.ResponseWriter, r *http.Request) {
	courtInput := model.Court{
		Name:           extractFieldFromForm("name", r),
		Url:            "court.Url",
		Adress:         extractFieldFromForm("adress", r),
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

func (server *BasketServer) getCourtHandler(w http.ResponseWriter, r *http.Request) {
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
