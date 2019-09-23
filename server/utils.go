package server

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/yousseffarkhani/court/model"
)

// TODO : Extract with datamodel
func extractFieldFromForm(fieldName string, r *http.Request) string {
	return strings.TrimSpace(r.FormValue(fieldName))
}

func extractEmptyFieldErrors(data interface{}) map[string]string {
	errors := make(map[string]string)
	value := reflect.ValueOf(data)
	typeOfData := value.Type()
	if value.Kind() == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			if value.Field(i).String() == "" {
				errors[typeOfData.Field(i).Name] = "Empty field"
			}
		}
	}
	return errors
}

func getInput(w http.ResponseWriter, r *http.Request) (model.User, map[string]string) {
	username := extractFieldFromForm("username", r)
	password := extractFieldFromForm("password", r)

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

func redirectToIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}
