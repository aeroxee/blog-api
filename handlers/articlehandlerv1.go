package handlers

import (
	"fmt"
	"net/http"

	"github.com/Aeroxee/blog-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"
)

type ArticleHandlerV1 struct{}

func NewArticleHandlerV1() ArticleHandlerV1 {
	return ArticleHandlerV1{}
}

func (ArticleHandlerV1) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status := getQueryString(ctx.Request, "status", "PUBLISHED")
		page := getQueryInteger(ctx.Request, "page", 1)
		size := getQueryInteger(ctx.Request, "size", 10)
		orderBy := getQueryString(ctx.Request, "order_by", "created_at")
		desc := getQueryBool(ctx.Request, "desc", true)
		categoryId := getQueryInteger(ctx.Request, "category_id", 0)
		userId := getQueryInteger(ctx.Request, "user_id", 0)
		q := getQueryString(ctx.Request, "q", "")

		offset := (page - 1) * size

		// check status is not PUBLISHED or DRAFTED
		if (status != "DRAFTED") && (status != "PUBLISHED") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Harap gunakan query status dengan nilai/value `DRAFTED` atau `PUBLISHED` saja.",
			})
			return
		}

		field := models.Article{
			CategoryID: categoryId,
			UserID:     userId,
			Status:     models.StatusArticle(status),
		}

		articleModel := models.NewArticleModel(models.GetDB())
		articles := articleModel.GetArticleWithFilter(field, offset, size, clause.OrderByColumn{
			Column: clause.Column{
				Name: orderBy,
			},
			Desc: desc,
		}, q)

		// count total of articles
		var count int64
		models.GetDB().Model(field).Where("status = ?", status).Count(&count)

		nextPage := page + 1
		var prevPage int

		if page > 1 {
			prevPage = page - 1
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":         "success",
			"total":          count,
			"page":           page,
			"size":           size,
			"order_by":       orderBy,
			"desc":           desc,
			"status_article": status,
			"category_id":    categoryId,
			"user_id":        userId,
			"next_page":      nextPage,
			"prev_page":      prevPage,
			"articles":       articles,
		})
	}
}

func (ArticleHandlerV1) Detail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")

		articleModel := models.NewArticleModel(models.GetDB())
		userModel := models.NewUserModel(models.GetDB())

		owner, err := userModel.GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "User tidak dapat ditemukan.",
			})
			return
		}

		article, err := articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Artikel tidak dapat ditemukan.",
			})
			return
		}

		// if article.Status == models.DRAFTED {
		// 	ctx.JSON(http.StatusNotFound, gin.H{
		// 		"status":  "error",
		// 		"message": "Artikel tidak dapat ditemukan.",
		// 	})
		// 	return
		// }

		// add 1 views
		if article.Status == models.PUBLISHED {
			article.Views += 1
			models.GetDB().Save(&article)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Berhasil mendapatkan data artikel.",
			"article": article,
		})
	}
}

func (ArticleHandlerV1) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		var payload ArticleCreatePayload
		err = ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
				"test":    err.Error(),
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
					Tag:   err.ActualTag(),
				})
			}

			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Validasi error",
				"errors":  validations,
			})
			return
		}

		article := models.Article{
			Title:      payload.Title,
			CategoryID: payload.CategoryID,
			Slug:       slug.MakeLang(payload.Title, "id"),
			Content:    payload.Content,
			UserID:     thisUser.ID,
			Status:     payload.Status,
		}

		articleModel := models.NewArticleModel(models.GetDB())
		err = articleModel.CreateArticle(&article)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Berhasil menambahkan satu artikel.",
			"article": article,
		})
	}
}

func (ArticleHandlerV1) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")

		articleModel := models.NewArticleModel(models.GetDB())
		userModel := models.NewUserModel(models.GetDB())

		owner, err := userModel.GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "User tidak dapat ditemukan.",
			})
			return
		}

		article, err := articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Artikel dengan slug `%s` tidak dapat ditemukan.", slugArticle),
			})
			return
		}

		// check thisUser is owner by this article
		if article.UserID != thisUser.ID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Kamu tidak mempunyai ijin untuk mengedit artikel milik orang lain.",
			})
			return
		}

		var payload ArticleUpdatePayload
		err = ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
			})
			return
		}

		if payload.Title != "" {
			article.Title = payload.Title
			article.Slug = slug.MakeLang(payload.Title, "id")
		}
		if payload.Content != "" {
			article.Content = payload.Content
		}
		if payload.Status != "" {
			article.Status = payload.Status
		}
		if payload.CategoryID != 0 {
			article.CategoryID = payload.CategoryID
		}

		// save
		err = models.GetDB().Save(&article).Error
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Berhasil memperbaharui artikel.",
			"article": article,
		})
	}
}

func (ArticleHandlerV1) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan.",
			})
			return
		}

		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")

		articleModel := models.NewArticleModel(models.GetDB())
		userModel := models.NewUserModel(models.GetDB())

		owner, err := userModel.GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "User tidak dapat ditemukan.",
			})
			return
		}

		article, err := articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Artikel dengan slug `%s` tidak dapat ditemukan.", slugArticle),
			})
			return
		}

		// check thisUser is owner by this article
		if article.UserID != thisUser.ID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Kamu tidak mempunyai ijin untuk menghapus artikel milik orang lain.",
			})
			return
		}

		// delete
		err = models.GetDB().Delete(&article).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}
