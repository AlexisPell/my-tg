package router

import (
	"github.com/alexispell/my-tg/users/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Auth routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/logout", handler.Logout)
	}

	// User routes
	userGroup := r.Group("/users")
	{
		userGroup.POST("/create", handler.CreateUser)
		userGroup.GET("/:id", handler.GetUserByID)
	}

	return r
}
