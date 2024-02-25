package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StatusArticle string

const (
	PUBLISHED StatusArticle = "PUBLISHED"
	DRAFTED   StatusArticle = "DRAFTED"
)

type Article struct {
	ID         int           `gorm:"primaryKey" json:"id"`
	UserID     int           `json:"user_id"`
	CategoryID int           `json:"category_id"`
	Title      string        `gorm:"size:255" json:"title"`
	Slug       string        `gorm:"size:255;unique" json:"slug"`
	Content    string        `gorm:"type:text" json:"content"`
	Status     StatusArticle `gorm:"size:9;default:DRAFTED" json:"status"`
	Views      int64         `gorm:"default:0" json:"views"`
	UpdatedAt  time.Time     `json:"updated_at"`
	CreatedAt  time.Time     `json:"created_at"`
	Comments   []Comment     `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
}

type ArticleModel struct {
	db *gorm.DB
}

func NewArticleModel(db *gorm.DB) ArticleModel {
	return ArticleModel{db}
}

func (a *ArticleModel) CreateArticle(article *Article) error {
	userModel := NewUserModel(GetDB())
	u, _ := userModel.GetUserByID(article.UserID)

	l := Log{
		UserID:      article.UserID,
		Title:       "Membuat artikel baru",
		Description: "Anda membuat artikel baru dengan nama " + article.Title,
		Url:         fmt.Sprintf("/v1/articles/%s/%s", u.Username, article.Slug),
	}

	if err := NewLog(&l); err != nil {
		fmt.Println(err)
	}
	return a.db.Create(article).Error
}

func (a *ArticleModel) GetArticleWithFilter(field Article, offset, limit int, order clause.OrderByColumn, q string) []Article {
	var articles []Article
	if q == "" {
		a.db.Model(&Article{}).Where(&field).
			Order(order).Preload("Comments").Offset(offset).Limit(limit).Find(&articles)
	} else {
		a.db.Model(&Article{}).Where(&field).Where("title LIKE ?", "%"+q+"%").
			Or("content LIKE ? ", "%"+q+"%").
			Order(order).Preload("Comments").Offset(offset).Limit(limit).Find(&articles)
	}
	return articles
}

func (a *ArticleModel) GetArticleBySlugAndUsername(slug string, userId int) (Article, error) {
	var article Article
	err := a.db.Model(&Article{}).Where("slug = ?", slug).Where("user_id = ?", userId).Preload("Comments").First(&article).Error
	return article, err
}
