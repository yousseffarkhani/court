package courtdb

import (
	"encoding/json"
	"fmt"
	"log"
	"model"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type CourtStore struct {
	db *gorm.DB
}

var (
	driverName = "postgres"
	host       = "db"
	port       = "5432"
	user       = os.Getenv("POSTGRES_USER")
	password   = os.Getenv("POSTGRES_PASSWORD")
	dbname     = "basket"
)

func NewCourtStore() (*CourtStore, error) {
	db, err := InitialMigration()
	if err != nil {
		return &CourtStore{}, err
	}
	courtStore := &CourtStore{db}
	courtStore.InitDB()
	return courtStore, nil
}

func (db *CourtStore) Close() error {
	return db.db.Close()
}

/* Database Initialisation */

func InitialMigration() (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable", host, port, user, password)
	err := Reset(driverName, psqlInfo, dbname)
	if err != nil {
		return nil, err
	}
	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err := gorm.Open(driverName, psqlInfo)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&model.Court{})
	db.AutoMigrate(&model.User{})
	return db, nil
}

func Reset(driverName, dataSource, dbname string) error {
	db, err := gorm.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	resetDB(db, dbname)

	return db.Close()
}

func resetDB(db *gorm.DB, name string) {
	db.Exec("DROP DATABASE IF EXISTS " + name)
	createDB(db, name)
}

func createDB(db *gorm.DB, name string) {
	db.Exec("CREATE DATABASE " + name)
}

func (db *CourtStore) InitDB() {
	file, err := os.Open("assets/courts.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var courts []model.Court
	err = decoder.Decode(&courts)
	if err != nil {
		log.Fatalln(err)
	}

	for _, court := range courts {
		db.AddCourt(court)
	}
}

/* Court CRUD methods */
func (db *CourtStore) AddCourt(newCourt model.Court) error {
	if _, err := db.GetCourtByName(newCourt.Name); err != nil {
		db.db.Create(&newCourt)
		fmt.Println("New court successfully created")
		return nil
	}
	return duplicateError{"This court already exists"}
}

func (db *CourtStore) GetAllCourts() []model.Court {
	var courts []model.Court
	db.db.Find(&courts)
	return courts
}

func (db *CourtStore) GetCourt(id int) model.Court {
	var court model.Court
	db.db.Where("id=?", id).Find(&court)
	return court
}

func (db *CourtStore) GetCourtByName(name string) (model.Court, error) {
	var court model.Court
	db.db.Where("name=?", name).Find(&court)
	if (court == model.Court{}) {
		return model.Court{}, duplicateError{"This court doesn't exist"}
	}
	return court, nil
}

/* User CRUD methods */
func (db *CourtStore) AddUser(newUser model.User) error {
	if _, err := db.GetUser(newUser.Username); err != nil {
		db.db.Create(&newUser)
		fmt.Println("New user successfully created")
		return nil
	}
	return duplicateError{"This username already exists"}
}

func (db *CourtStore) GetUser(username string) (model.User, error) {
	var user model.User
	db.db.Where("username=?", username).Find(&user)
	if (user == model.User{}) {
		return model.User{}, duplicateError{"This username doesn't exist"}
	}
	return user, nil
}

type duplicateError struct {
	description string
}

func (err duplicateError) Error() string {
	return err.description
}
