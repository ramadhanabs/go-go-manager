package v1

import (
	"database/sql"
	"go-go-manager/models"
	"go-go-manager/repositories"
	"go-go-manager/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	Repo *repositories.EmployeeRepository
}

func NewEmployeeHandler(db *sql.DB) *EmployeeHandler {
	return &EmployeeHandler{
		Repo: repositories.NewEmployeeRepository(db),
	}
}

func (h *EmployeeHandler) CreateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate the token
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing request token"})
			return
		}

		if c.GetHeader("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing content-type"})
			return
		}

		auth = auth[7:] // Remove "Bearer " prefix
		v, err := utils.ValidateJWT(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Proceed with the handler logic
		var employee models.Employee
		if err := c.ShouldBindJSON(&employee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate gender
		if employee.Gender != models.Male && employee.Gender != models.Female {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender value"})
			return
		}

		// Check for duplicate identity number
		existingEmployee, err := h.Repo.GetEmployeeByIdentityNumber(employee.IdentityNumber)
		if err == nil && existingEmployee != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Identity number conflict"})
			return
		}

		// Check department id available or not
		_, err = models.FindDepartmentById(v.UserID, employee.DepartmentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Department ID"})
			return
		}

		// Add employee to the database
		if err := h.Repo.AddEmployee(employee); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"departmentId":     employee.DepartmentID,
			"name":             employee.Name,
			"identityNumber":   employee.IdentityNumber,
			"gender":           employee.Gender,
			"employeeImageUri": employee.EmployeeImageURI,
		})
	}
}

type EmployeeResponse struct {
	IdentityNumber   string        `json:"identityNumber"`
	Name             string        `json:"name"`
	Gender           models.Gender `json:"gender"`
	DepartmentID     string        `json:"departmentId"`
	EmployeeImageURI string        `json:"employeeImageUri"`
}

func (h *EmployeeHandler) GetEmployees() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate the token
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing request token"})
			return
		}

		auth = auth[7:] // Remove "Bearer " prefix
		_, err := utils.ValidateJWT(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Proceed with the handler logic
		filters := make(map[string]string)

		// Extract query parameters
		if identityNumber := c.Query("identityNumber"); identityNumber != "" {
			filters["identityNumber"] = identityNumber
		}
		if name := c.Query("name"); name != "" {
			filters["name"] = name
		}
		if gender := c.Query("gender"); gender != "" {
			filters["gender"] = gender
		}
		if departmentID := c.Query("departmentId"); departmentID != "" {
			filters["departmentId"] = departmentID
		}

		// Pagination
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		filters["limit"] = strconv.Itoa(limit)
		filters["offset"] = strconv.Itoa(offset)

		// Fetch employees from the database
		employees, err := h.Repo.FilterEmployees(filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
			return
		}

		response := make([]EmployeeResponse, 0)
		for _, employee := range employees {
			response = append(response, EmployeeResponse{
				IdentityNumber:   employee.IdentityNumber,
				Name:             employee.Name,
				Gender:           employee.Gender,
				DepartmentID:     employee.DepartmentID,
				EmployeeImageURI: employee.EmployeeImageURI,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h *EmployeeHandler) UpdateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate the token
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing request token"})
			return
		}

		if c.GetHeader("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing content-type"})
			return
		}

		auth = auth[7:] // Remove "Bearer " prefix
		_, err := utils.ValidateJWT(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		identityNumber := c.Param("identityNumber")

		// First check if employee exists
		existingEmployee, err := h.Repo.GetEmployeeByIdentityNumber(identityNumber)
		if err != nil || existingEmployee == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		var updatedEmployee models.Employee
		if err := c.ShouldBindJSON(&updatedEmployee); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate gender
		if updatedEmployee.Gender != models.Male && updatedEmployee.Gender != models.Female {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender value"})
			return
		}

		// Update employee in the database
		if err := h.Repo.UpdateEmployee(identityNumber, updatedEmployee); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}

		c.JSON(http.StatusOK, updatedEmployee)
	}
}

func (h *EmployeeHandler) DeleteEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate the token
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing request token"})
			return
		}

		auth = auth[7:] // Remove "Bearer " prefix
		_, err := utils.ValidateJWT(auth)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Proceed with the handler logic
		identityNumber := c.Param("identityNumber")

		_, err = h.Repo.GetEmployeeByIdentityNumber(identityNumber)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
			return
		}

		// Delete employee from the database
		if err := h.Repo.DeleteEmployee(identityNumber); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Employee deleted"})
	}
}
