package routes

import (
	"net/http"

	"github.com/doxanocap/reactNative/dino-back/pkg/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()
	api := r.Group("/api")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PATCH", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Accept-Encoding", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))
	api.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "ping"})
	})
	api.POST("/signUp", controllers.SignUp)
	api.POST("/signIn", controllers.SignIn)
	api.GET("/signOut", controllers.SignOut)
	api.GET("/get-user-info", controllers.User)
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
