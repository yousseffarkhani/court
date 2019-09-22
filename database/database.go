package database

import (
	"fmt"
	"os"

	"github.com/yousseffarkhani/court/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

// TODO : Rename to Database
type CourtStore struct {
	*gorm.DB
}

var (
	driverName = "postgres"
	host       = "db"
	port       = "5432"
	user       = os.Getenv("POSTGRES_USER")
	password   = os.Getenv("POSTGRES_PASSWORD")
	dbname     = "basket"
)

func NewCourtStore(file *os.File) (*CourtStore, error) {
	defer file.Close()

	db, err := initializeDB()
	if err != nil {
		return nil, fmt.Errorf("Problem initializing DB, %v", err)
	}

	courtStore := &CourtStore{db}

	courtStore.initData(file)

	return courtStore, nil
}

func (db *CourtStore) initData(file *os.File) error {
	courts, err := model.NewCourts(file)
	if err != nil {
		return fmt.Errorf("Problem loading courts from file %s, %v", file.Name(), err)
	}

	for _, court := range courts {
		db.AddCourt(court)
	}
	return nil
}

/* Database Initialisation */

func initializeDB() (*gorm.DB, error) {
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
	db.AutoMigrate(&model.Comment{})

	return db, nil
}

func (db *CourtStore) Close() error {
	return db.Close()
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

/* Utils */
type customError struct {
	description string
}

func (err customError) Error() string {
	return err.description
}
