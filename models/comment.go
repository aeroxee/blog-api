package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	ArticleID int       `json:"article_id"`
	UserID    int       `json:"user_id"`
	Text      string    `gorm:"type:text" json:"text"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentModel struct {
	db *gorm.DB
}

func NewCommentModel(db *gorm.DB) CommentModel {
	return CommentModel{db: db}
}

func (c *CommentModel) CreateComment(comment *Comment) error {
	return c.db.Create(comment).Error
}

func (c *CommentModel) GetCommentByID(id int) (Comment, error) {
	var comment Comment
	err := c.db.Model(&Comment{}).Where("id = ?", id).First(&comment).Error
	return comment, err
}

func (c *CommentModel) GetCommentByUserID(userId int) []Comment {
	var comments []Comment
	c.db.Model(&Comment{}).Where("user_id = ?", userId).Find(&comments)
	return comments
}

func (c *CommentModel) GetCommentByArticleID(articleId int) []Comment {
	var comments []Comment
	c.db.Model(&Comment{}).Where("article_id = ?", articleId).Find(&comments)
	return comments
}
