package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Aeroxee/blog-api/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CommentHandlerV1 struct{}

func NewCommentHandlerV1() CommentHandlerV1 {
	return CommentHandlerV1{}
}

func (CommentHandlerV1) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")
		owner, err := models.NewUserModel(models.GetDB()).GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		articleModel := models.NewArticleModel(models.GetDB())
		article, err := articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Artikel yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		commentModel := models.NewCommentModel(models.GetDB())
		comments := commentModel.GetCommentByArticleID(article.ID)
		ctx.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"message":  "",
			"comments": comments,
		})
	}
}

func (CommentHandlerV1) Detail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")
		owner, err := models.NewUserModel(models.GetDB()).GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		articleModel := models.NewArticleModel(models.GetDB())
		_, err = articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Artikel yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		// get comment id by param
		commentId := ctx.Param("commentId")
		commentIdInt, err := strconv.Atoi(commentId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		commentModel := models.NewCommentModel(models.GetDB())
		comment, err := commentModel.GetCommentByID(commentIdInt)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "",
			"comment": comment,
		})
	}
}

func (CommentHandlerV1) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan",
			})
			return
		}

		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")
		owner, err := models.NewUserModel(models.GetDB()).GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		articleModel := models.NewArticleModel(models.GetDB())
		article, err := articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Artikel yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		var payload CommentPayload
		err = ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
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
				"message": "Validasi field error",
				"errors":  validations,
			})
			return
		}

		comment := models.Comment{
			UserID:    thisUser.ID,
			ArticleID: article.ID,
			Text:      payload.Text,
		}

		// save
		commentModel := models.NewCommentModel(models.GetDB())
		err = commentModel.CreateComment(&comment)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"comment": comment,
			"message": "Berhasil mengomentari artikel.",
		})
	}
}

func (CommentHandlerV1) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan",
			})
			return
		}

		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")
		owner, err := models.NewUserModel(models.GetDB()).GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		articleModel := models.NewArticleModel(models.GetDB())
		_, err = articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Artikel yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		// get comment id by param
		commentId := ctx.Param("commentId")
		commentIdInt, err := strconv.Atoi(commentId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		commentModel := models.NewCommentModel(models.GetDB())
		comment, err := commentModel.GetCommentByID(commentIdInt)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		payload := struct {
			Text string `json:"text"`
		}{}
		err = ctx.ShouldBindJSON(&payload)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Payload yang anda gunakan bukan bertipe json. Harap kirim data menggunakan json saja.",
			})
			return
		}

		// check thisUser is owner of this comment
		if thisUser.ID != comment.UserID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Anda tidak dapat memperbaharui komentar milik orang lain.",
			})
			return
		}

		if payload.Text != "" {
			comment.Text = payload.Text
		}

		// save
		models.GetDB().Save(&comment)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Berhasil memperbaharui komentar.",
			"comment": comment,
		})
	}
}

func (CommentHandlerV1) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		thisUser, err := getUserContext(ctx.Request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Autentikasi dibutuhkan",
			})
			return
		}

		slugArticle := ctx.Param("slug")
		username := ctx.Param("username")
		owner, err := models.NewUserModel(models.GetDB()).GetUserUsername(username)
		if err != nil {
			ctx.JSON(http.StatusNotFound, nil)
			return
		}

		articleModel := models.NewArticleModel(models.GetDB())
		_, err = articleModel.GetArticleBySlugAndUsername(slugArticle, owner.ID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Artikel yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		// get comment id by param
		commentId := ctx.Param("commentId")
		commentIdInt, err := strconv.Atoi(commentId)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		commentModel := models.NewCommentModel(models.GetDB())
		comment, err := commentModel.GetCommentByID(commentIdInt)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Halaman yang anda tuju tidak dapat ditemukan.",
			})
			return
		}

		// check thisUser is owner of this comment
		if thisUser.ID != comment.UserID {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Anda tidak dapat menghapus komentar milik orang lain.",
			})
			return
		}

		// save
		models.GetDB().Delete(&comment)

		ctx.JSON(http.StatusNoContent, nil)
	}
}
