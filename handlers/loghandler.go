package handlers

import (
	"net/http"
	"strconv"

	"github.com/Aeroxee/blog-api/models"
	"github.com/gin-gonic/gin"
)

type LogHandlerV1 struct{}

func NewLogHandlerV1() LogHandlerV1 {
	return LogHandlerV1{}
}

func (LogHandlerV1) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		logs := models.GetAllLogFromUserID(userIdInt)
		var count int64
		models.GetDB().Model(&models.Log{}).Count(&count)

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "",
			"logs":    logs,
			"total":   count,
		})
	}
}

// TODO: Detail, Post, Put, Delete
