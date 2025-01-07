package v1

import (
	"go-go-manager/models"
	"go-go-manager/utils"
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

type UpdateDepartmentRequest struct {
	Name string `json:"name" binding:"required,min=4,max=33"`
}

func UpdateDepartment(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	auth = auth[7:]

	v, err := utils.ValidateJWT(auth)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	var req UpdateDepartmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	departmentId := c.Param("departmentId")
	if len(departmentId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "departmentId is required"})
		return
	}

	_, err = models.FindDepartmentById(v.UserID, departmentId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "departmentId is not found"})
		return

	} else {
		department, err := models.UpdateDepartment(departmentId, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"departmentId": department.ID,
			"name":         department.Name,
		})
	}
}

func DeleteDepartment(c *gin.Context) {

	auth := c.GetHeader("Authorization")
	auth = auth[7:]

	v, err := utils.ValidateJWT(auth)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	departmentId := c.Param("departmentId")
	if len(departmentId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "departmentId is required"})
		return
	}

	_, err = models.FindDepartmentById(v.UserID, departmentId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "department not found"})
		return
	}

	employeeCount, err := models.CountEmployeesByDepartment(departmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check for associated employees"})
		return
	}

	if employeeCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Still contain employee"})
		return
	}

	err = models.DeleteDepartment(departmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Department deleted")
}
