package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func UploadHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		uid := uuid.NewString()
		randomString := strings.ReplaceAll(uid, "-", "")

		h, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ext := filepath.Ext(h.Filename)
		filename := randomString + ext

		destination := fmt.Sprintf("media/uploads/%s", filename)

		// save
		if err = ctx.SaveUploadedFile(h, destination); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Upload file berhasil",
			"url":     fmt.Sprintf("http://%s:%s/%s", os.Getenv("HOSTNAME"), os.Getenv("PORT"), destination),
		})
	}
}
