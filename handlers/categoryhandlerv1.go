package handlers

import (
	"fmt"
	"net/http"

	"github.com/Aeroxee/blog-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
)

type CategoryHandlerV1 struct{}

func NewCategoryHandlerV1() CategoryHandlerV1 {
	return CategoryHandlerV1{}
}

func (CategoryHandlerV1) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		articleModel := models.NewArticleModel(models.GetDB())
		id := getQueryInteger(ctx.Request, "id", 0)
		if id != 0 {
			category, err := articleModel.GetCategoryByID(id)
			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": "Kategori tidak dapat ditemukan.",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"status":   "success",
				"message":  "",
				"category": category,
			})
			return
		}

		categories := articleModel.GetAllCategory()

		ctx.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"categories": categories,
		})
	}
}

func (CategoryHandlerV1) Detail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slugCategory := ctx.Param("slug")

		articleModel := models.NewArticleModel(models.GetDB())
		category, err := articleModel.GetCategoryBySlug(slugCategory)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"message":  "",
			"category": category,
		})
	}
}

func (CategoryHandlerV1) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		if !thisUser.IsAdmin {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Anda tidak diijinkan untuk mengakses metode ini.",
			})
			return
		}

		var payload CategoryCreatePayload
		err = ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		validate = validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(payload)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				fmt.Println(err.Error())
				return
			}

			var validations []Validation
			for _, err := range err.(validator.ValidationErrors) {
				validations = append(validations, Validation{
					Field: err.Field(),
					Tag:   err.Tag(),
				})
			}

			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi error",
				"errors":  validations,
			})
			return
		}

		category := models.Category{
			Title:       payload.Title,
			Slug:        slug.MakeLang(payload.Title, "id"),
			Description: payload.Description,
		}

		articleModel := models.NewArticleModel(models.GetDB())

		err = articleModel.CreateCategory(&category)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":   "success",
			"message":  "Berhasil menambahkan kategori baru.",
			"category": category,
		})
	}
}

func (CategoryHandlerV1) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		if !thisUser.IsAdmin {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Anda tidak diijinkan untuk mengakses metode ini.",
			})
			return
		}

		slugCategory := ctx.Param("slug")
		articleModel := models.NewArticleModel(models.GetDB())

		category, err := articleModel.GetCategoryBySlug(slugCategory)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kategory tidak dapat ditemukan.",
			})
			return
		}

		payload := struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}{}
		err = ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		if payload.Title != "" {
			category.Title = payload.Title
			category.Slug = slug.MakeLang(payload.Title, "id")
		}
		if payload.Description != "" {
			category.Description = payload.Description
		}

		models.GetDB().Save(&category)
		ctx.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"message":  "Berhasil memperbaharui kategory.",
			"category": category,
		})
	}
}

func (CategoryHandlerV1) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		if !thisUser.IsAdmin {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Anda tidak diijinkan untuk mengakses metode ini.",
			})
			return
		}

		slugCategory := ctx.Param("slug")
		articleModel := models.NewArticleModel(models.GetDB())

		category, err := articleModel.GetCategoryBySlug(slugCategory)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Kategory tidak dapat ditemukan.",
			})
			return
		}

		models.GetDB().Unscoped().Delete(&category)

		ctx.JSON(http.StatusNoContent, nil)
	}
}
