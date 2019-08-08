package model

import (
	"github.com/jinzhu/gorm"
)

type Court struct {
	gorm.Model
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

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"primary_key"`
	Password string `json:"password"`
}

type TemplateData struct {
	Ressource  interface{}
	Errors     map[string]string
	ActionDone bool
}

type Comment struct {
	gorm.Model
	Message string `json:"message"`
	CourtID int    `json:"courtid"`
	// TODO : Author
	// Date
}
