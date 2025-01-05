package routes

import (
	v1 "go-go-manager/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1Group := router.Group("/api/v1")
	{
		v1Group.GET("/users", v1.GetUsers)
	}

	return router
}
