package courtdb

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
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

type Store interface {
	AddCourt(court Court) error
	GetCourts() ([]*Court, error)
	GetCourt() Court
}

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

/* Initialize database */

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
	db.AutoMigrate(&Court{})
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
	var courts []Court
	err = decoder.Decode(&courts)
	if err != nil {
		log.Fatalln(err)
	}

	for _, court := range courts {
		db.AddCourt(court)
	}
}

/* CRUD methods */
func Open() (*CourtStore, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password)
	db, err := gorm.Open(driverName, psqlInfo)
	if err != nil {
		return nil, err
	}
	return &CourtStore{db}, nil
}

func (db *CourtStore) AddCourt(court Court) {
	db.db.Create(&Court{
		Name:           court.Name,
		Url:            court.Url,
		Adress:         court.Adress,
		Arrondissement: court.Arrondissement,
		Longitude:      court.Longitude,
		Lattitude:      court.Lattitude,
		Dimensions:     court.Dimensions,
		Revetement:     court.Revetement,
		Decouvert:      court.Decouvert,
		Eclairage:      court.Eclairage,
	})
	fmt.Println("New court successfully created")
}

func (db *CourtStore) GetAllCourts() []Court {
	var courts []Court
	db.db.Find(&courts)
	return courts
}

func (db *CourtStore) GetCourt(id string) Court {
	var court Court
	db.db.Where("id=?", id).Find(&court)
	return court
}
