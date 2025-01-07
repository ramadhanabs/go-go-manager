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
		v1Group.GET("/user", v1.GetUsers)
		v1Group.PATCH("/user", v1.UpdateUser)
		v1Group.POST("/department", v1.CreateDepartment)
		v1Group.GET("/department", v1.GetDepartments)
		v1Group.PATCH("/department/:departmentId", v1.UpdateDepartment)
		v1Group.DELETE("/department/:departmentId", v1.DeleteDepartment)

	}

	return router
}
