package handlers

import (
	"github.com/Aeroxee/blog-api/models"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type RegisterPayload struct {
	FirstName string `json:"first_name" validate:"required,max=50"`
	LastName  string `json:"last_name" validate:"required,max=50"`
	Username  string `json:"username" validate:"required,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=5,max=15"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ArticleCreatePayload struct {
	Title      string               `json:"title" validate:"required"`
	Content    string               `json:"content" validate:"required"`
	CategoryID int                  `json:"category_id" validate:"required"`
	Status     models.StatusArticle `json:"status" validate:"required"`
}

type ArticleUpdatePayload struct {
	Title      string               `json:"title" `
	Content    string               `json:"content"`
	CategoryID int                  `json:"category_id"`
	Status     models.StatusArticle `json:"status"`
}

type CommentPayload struct {
	Text string `json:"text" validate:"required"`
}

type CategoryCreatePayload struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
