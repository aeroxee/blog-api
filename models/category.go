package models

import (
	"time"

	"gorm.io/gorm/clause"
)

type Category struct {
	ID          int        `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:50;unique" json:"title"`
	Slug        string     `gorm:"size:60;unique" json:"slug"`
	Description string     `gorm:"type:text" json:"description"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	Articles    []*Article `gorm:"foreignKey:CategoryID" json:"articles"`
}

func (a *ArticleModel) CreateCategory(category *Category) error {
	return a.db.Create(category).Error
}

func (a *ArticleModel) GetAllCategory() []Category {
	var categories []Category
	a.db.Model(&Category{}).Preload("Articles").Order(clause.OrderByColumn{
		Column: clause.Column{Name: "updated_at"},
		Desc:   true,
	}).Find(&categories)
	return categories
}

func (a *ArticleModel) GetCategoryByID(id int) (Category, error) {
	var category Category
	err := a.db.Model(&Category{}).Where("id = ?", id).Preload("Articles").First(&category).Error
	return category, err
}

func (a *ArticleModel) GetCategoryBySlug(slug string) (Category, error) {
	var category Category
	err := a.db.Model(&Category{}).Where("slug = ?", slug).Preload("Articles").First(&category).Error
	return category, err
}
