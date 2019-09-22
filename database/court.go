package database

import (
	"fmt"

	"github.com/yousseffarkhani/court/model"
)

/* Court CRUD methods */
func (db *CourtStore) AddCourt(newCourt model.Court) error {
	if _, err := db.GetCourtByName(newCourt.Name); err != nil {
		db.Create(&newCourt)
		fmt.Println("New court successfully created")
		return nil
	}
	return customError{"This court already exists"}
}

func (db *CourtStore) GetAllCourts() []model.Court {
	var courts []model.Court
	db.Find(&courts)
	return courts
}

func (db *CourtStore) GetCourt(id int) model.Court {
	var court model.Court
	db.Where("id=?", id).Find(&court)
	return court
}

func (db *CourtStore) GetCourtByName(name string) (model.Court, error) {
	var court model.Court
	db.Where("name=?", name).Find(&court)
	if (court == model.Court{}) {
		return model.Court{}, customError{"This court doesn't exist"}
	}
	return court, nil
}
