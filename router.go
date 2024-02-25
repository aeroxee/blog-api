package blogapi

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Static("/media", "./media")

	c := cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	})
	r.Use(c)

	// group v1
	v1 := r.Group("/v1")
	mainControllerNoAuth(v1)

	// user group
	userGroupV1 := v1.Group("/user")
	userGroupV1.Use(authMiddleware())
	userController(userGroupV1)

	// article group v1 no auth
	articleGroupV1NoAuth := v1.Group("/articles")
	articleControllerNoAuth(articleGroupV1NoAuth)

	// article group v1 with auth
	articleGroupV1WithAuth := v1.Group("/articles")
	articleGroupV1WithAuth.Use(authMiddleware())
	articleControllerWithAuth(articleGroupV1WithAuth)

	// category group v1 no auth
	categoryGroupV1NoAuth := v1.Group("/categories")
	categoryControllerNoAuth(categoryGroupV1NoAuth)

	// category group v1 with auth
	categoryGroupV1WithAuth := v1.Group("/categories")
	categoryGroupV1WithAuth.Use(authMiddleware())
	categoryControllerWithAuth(categoryGroupV1WithAuth)

	// log handler
	logGroup := v1.Group("/logs")
	logGroup.Use(authMiddleware())
	logController(logGroup)

	// r.Run("192.168.116.225:8000")
	r.Run(":8000")
}
