package v1

import (
	"go-go-manager/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DepartmentRequest struct {
	Name string `json:"name" binding:"required,min=4,max=33"`
}

func CreateDepartment(c *gin.Context) {
	var req DepartmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := models.FindDepartmentByName(req.Name)
	if err != nil {
		department, err := models.CreateDepartment(req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"departmentId": department.ID,
			"name":         department.Name,
		})
	} else {
		c.JSON(http.StatusConflict, gin.H{"error": "Department already exist"})
		return
	}

}
