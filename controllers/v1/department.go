package v1

import (
	"go-go-manager/models"
	"go-go-manager/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type DepartmentRequest struct {
	Name string `json:"name" binding:"required,min=4,max=33"`
}

func CreateDepartment(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return
	}

	auth = auth[7:]
	v, err := utils.ValidateJWT(auth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req DepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department name cannot be empty"})
		return
	}

	_, err = models.FindDepartmentByName(req.Name)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Department already exists"})
		return
	}

	department, err := models.CreateDepartment(req.Name, v.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"departmentId": department.ID,
		"name":         department.Name,
	})
}

func GetDepartments(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return
	}

	auth = auth[7:]
	v, err := utils.ValidateJWT(auth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	limit := 5
	offset := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	name := c.Query("name")

	departments, err := models.GetDepartments(v.UserID, limit, offset, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]gin.H, 0)
	for _, dept := range departments {
		response = append(response, gin.H{
			"departmentId": strconv.Itoa(int(dept.ID)),
			"name":         dept.Name,
		})
	}

	if len(response) == 0 {
		response = []gin.H{}
	}

	c.JSON(http.StatusOK, response)
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
