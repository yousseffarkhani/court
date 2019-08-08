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
	db.AutoMigrate(&model.Comment{})
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
	return customError{"This court already exists"}
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
		return model.Court{}, customError{"This court doesn't exist"}
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
	return customError{"This username already exists"}
}

func (db *CourtStore) GetUser(username string) (model.User, error) {
	var user model.User
	db.db.Where("username=?", username).Find(&user)
	if (user == model.User{}) {
		return model.User{}, customError{"This username doesn't exist"}
	}
	return user, nil
}

/* Comment CRUD methods */
func (db *CourtStore) AddComment(newComment model.Comment) error {
	if err := db.db.Create(&newComment).Error; err != nil {
		return customError{"Couldn't add comment."}
	}
	fmt.Println("New comment added successfully")
	return nil
}

func (db *CourtStore) GetComments(id int) ([]model.Comment, error) {
	var comments []model.Comment
	err := db.db.Where("court_id=?", id).Find(&comments).Error
	if err != nil {
		return comments, customError{"Couldn't retrieve comments."}
	}
	return comments, nil
}

func (db *CourtStore) GetComment(commentId int) (model.Comment, error) {
	var comment model.Comment
	err := db.db.Where("id=?", commentId).Find(&comment).Error
	if gorm.IsRecordNotFoundError(err) {
		return comment, customError{"Comment doesn't exist."}
	} else if err != nil {
		return comment, customError{"Couldn't retrieve comment."}
	}
	return comment, nil
}

func (db *CourtStore) DeleteComment(id uint) error {
	comment, err := db.GetComment(int(id))
	if err != nil {
		return err
	} //TODO : Add court_id as criteria
	if err := db.db.Unscoped().Where("id=?", id).Delete(&comment).Error; err != nil {
		return customError{"Couldn't delete comment."}
	}
	fmt.Println("Comment successfully deleted")
	return nil
}

func (db *CourtStore) UpdateComment(updatedMessage model.Comment) error {
	comment, err := db.GetComment(int(updatedMessage.ID))
	if err != nil {
		return err
	}
	comment.Message = updatedMessage.Message
	if err := db.db.Save(&comment).Error; err != nil {
		return customError{"Couldn't update comment."}
	}
	fmt.Println("Comment successfully updated")
	return nil
}

/* Utils */
type customError struct {
	description string
}

func (err customError) Error() string {
	return err.description
}
