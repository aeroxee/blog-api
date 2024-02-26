package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	FirstName  string    `gorm:"size:50" json:"first_name"`
	LastName   string    `gorm:"size:50" json:"last_name"`
	Username   string    `gorm:"size:50;unique" json:"username"`
	Email      string    `gorm:"size:50;unique" json:"email"`
	Avatar     *string   `gorm:"size:255" json:"avatar"`
	Password   string    `gorm:"size:128" json:"-"`
	IsActive   bool      `gorm:"default:false" json:"is_active"`
	IsAdmin    bool      `gorm:"default:false" json:"is_admin"`
	UpdatedAt  time.Time `json:"updated_at"`
	DateJoined time.Time `gorm:"autoCreateTime" json:"date_joined"`
	Articles   []Article `gorm:"foreignKey:UserID" json:"articles"`
	Comments   []Comment `gorm:"foreignKey:UserID" json:"comments"`
}

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db}
}

func (u *UserModel) CreateUser(user *User) error {
	return u.db.Create(user).Error
}

func (u *UserModel) GetUserByID(id int) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("id = ?", id).Preload("Articles").Preload("Comments").First(&user).Error
	return user, err
}

func (u *UserModel) GetUserEmail(email string) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("email = ?", email).Preload("Articles").Preload("Comments").First(&user).Error
	return user, err
}

func (u *UserModel) GetUserUsername(username string) (User, error) {
	var user User
	err := u.db.Model(&User{}).Where("username = ?", username).Preload("Articles").Preload("Comments").First(&user).Error
	return user, err
}
