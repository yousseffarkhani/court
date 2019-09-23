package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/yousseffarkhani/court/model"
)

/* Comment CRUD methods */
func (db *CourtStore) AddComment(newComment model.Comment) error {
	if err := db.Create(&newComment).Error; err != nil {
		return customError{"Couldn't add comment."}
	}
	fmt.Println("New comment added successfully")
	return nil
}

func (db *CourtStore) GetComments(id int) ([]model.Comment, error) {
	var comments []model.Comment
	err := db.Where("court_id=?", id).Find(&comments).Error
	if err != nil {
		return comments, customError{"Couldn't retrieve comments."}
	}
	return comments, nil
}

func (db *CourtStore) GetComment(courtID, commentID int) (model.Comment, error) {
	var comment model.Comment
	err := db.Where("court_id=? AND id=? ", courtID, commentID).Find(&comment).Error
	if gorm.IsRecordNotFoundError(err) {
		return comment, customError{"Comment doesn't exist."}
	} else if err != nil {
		return comment, customError{"Couldn't retrieve comment."}
	}
	return comment, nil
}

func (db *CourtStore) DeleteComment(courtID, commentID int) error {
	if err := db.Unscoped().Where("court_id=? AND id=? ", courtID, commentID).Delete(&model.Comment{}).Error; err != nil {
		return customError{"Couldn't delete comment."}
	}
	fmt.Println("Comment successfully deleted")
	return nil
}

func (db *CourtStore) UpdateComment(courtID, commentID int, message string) error {
	comment, err := db.GetComment(courtID, commentID)
	if err != nil {
		return err
	}
	comment.Message = message
	if err := db.Save(&comment).Error; err != nil {
		return customError{"Couldn't update comment."}
	}
	fmt.Println("Comment successfully updated")
	return nil
}
