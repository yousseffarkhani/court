package database

import (
	"fmt"

	"github.com/yousseffarkhani/court/model"
)

/* User CRUD methods */
func (db *CourtStore) AddUser(newUser model.User) error {
	if _, err := db.GetUser(newUser.Username); err != nil {
		db.Create(&newUser)
		fmt.Println("New user successfully created")
		return nil
	}
	return customError{"This username already exists"}
}

func (db *CourtStore) GetUser(username string) (model.User, error) {
	var user model.User
	db.Where("username=?", username).Find(&user)
	if (user == model.User{}) {
		return model.User{}, customError{"This username doesn't exist"}
	}
	return user, nil
}
