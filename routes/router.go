package routes

import (
	"database/sql"
	v1 "go-go-manager/controllers/v1"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	employeeHandler := v1.NewEmployeeHandler(db)

	v1Group := router.Group("/api/v1")
	{
		v1Group.POST("/auth", v1.AuthHandler)
		v1Group.GET("/user", v1.GetUsers)
		v1Group.PATCH("/user", v1.UpdateUser)
		v1Group.POST("/department", v1.CreateDepartment)
		v1Group.PATCH("/department/:departmentId", v1.UpdateDepartment)
		v1Group.DELETE("/department/:departmentId", v1.DeleteDepartment)

		// Employee routes
		v1Group.POST("/employee", employeeHandler.CreateEmployee())
		v1Group.GET("/employee", employeeHandler.GetEmployees())
		v1Group.PATCH("/employee/:identityNumber", employeeHandler.UpdateEmployee())
		v1Group.DELETE("/employee/:identityNumber", employeeHandler.DeleteEmployee())

	}

	return router
}
