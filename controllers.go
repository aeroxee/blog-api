package blogapi

import (
	"github.com/Aeroxee/blog-api/handlers"
	"github.com/gin-gonic/gin"
)

func mainControllerNoAuth(group *gin.RouterGroup) {
	userHandlerV1 := handlers.NewUserHandlerV1()
	group.POST("/register", userHandlerV1.Register())
	group.POST("/get-token", userHandlerV1.GetToken())

	// check email and username
	group.POST("/check-email", userHandlerV1.CheckEmail())
	group.POST("/check-username", userHandlerV1.CheckUsername())
	group.GET("/activate/:activationCode", userHandlerV1.ActivationHandler())

	group.GET("/get-user-from-username/:username", userHandlerV1.GetUserFromUsername())
	group.GET("/get-user-from-id/:userId", userHandlerV1.GetUserFromID())

	// upload
	group.POST("/upload", handlers.UploadHandler())
}

func userController(group *gin.RouterGroup) {
	userHandlerV1 := handlers.NewUserHandlerV1()

	group.GET("/auth", userHandlerV1.Auth())
}

func articleControllerNoAuth(group *gin.RouterGroup) {
	articleHandlerV1 := handlers.NewArticleHandlerV1()
	commentHandlerV1 := handlers.NewCommentHandlerV1()

	group.GET("", articleHandlerV1.Get())
	group.GET("/:username/:slug", articleHandlerV1.Detail())

	// comment
	group.GET("/:username/:slug/comment", commentHandlerV1.Get())
	group.GET("/:username/:slug/comment/:commentId", commentHandlerV1.Detail())
}

func articleControllerWithAuth(group *gin.RouterGroup) {
	articleHandlerV1 := handlers.NewArticleHandlerV1()
	commentHandlerV1 := handlers.NewCommentHandlerV1()

	group.POST("", articleHandlerV1.Create())
	group.PUT("/:username/:slug", articleHandlerV1.Update())
	group.DELETE("/:username/:slug", articleHandlerV1.Delete())

	// comment
	group.POST("/:username/:slug/comment", commentHandlerV1.Create())
	group.PUT("/:username/:slug/comment/:commentId", commentHandlerV1.Update())
	group.DELETE("/:username/:slug/comment/:commentId", commentHandlerV1.Delete())
}

func categoryControllerNoAuth(group *gin.RouterGroup) {
	categoryHandlerV1 := handlers.NewCategoryHandlerV1()

	group.GET("", categoryHandlerV1.Get())
	group.GET("/:slug", categoryHandlerV1.Detail())
}

func categoryControllerWithAuth(group *gin.RouterGroup) {
	categoryHandlerV1 := handlers.NewCategoryHandlerV1()

	group.POST("", categoryHandlerV1.Create())
	group.PUT("/:slug", categoryHandlerV1.Update())
	group.DELETE("/:slug", categoryHandlerV1.Delete())
}

func logController(group *gin.RouterGroup) {
	logHandlerV1 := handlers.NewLogHandlerV1()
	group.GET("/:userId", logHandlerV1.Get())
}
