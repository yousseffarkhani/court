package server

import (
	"net/http"
	"reflect"
	"strings"
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
