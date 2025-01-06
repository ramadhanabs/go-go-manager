package routes

import (
	v1 "go-go-manager/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1Group := router.Group("/api/v1")
	{
		v1Group.POST("/auth", v1.AuthHandler)
		v1Group.POST("/department", v1.CreateDepartment)
	}

	return router
}
