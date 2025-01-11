package main

import (
	"fmt"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const PORT = "http://10.0.7.99"
const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJleHAiOjE3MzY2NzM2NTksImlhdCI6MTczNjU4NzI1OX0.6LQI_-1agUDnMaO-tvLtx21uCobZ3k0QcXjok1ZAEL8"

func TestSignupAPI(t *testing.T) {
	// Create a new httpexpect instance
	e := httpexpect.New(t, PORT)

	e.POST("/api/v1/auth").
		WithJSON(map[string]string{
			"email":    "test@test.com",
			"password": "password",
			"action":   "signup",
		}).
		Expect().
		Status(409)
}

func TestLoginAPI(t *testing.T) {
	// Create a new httpexpect instance
	e := httpexpect.New(t, PORT)

	e.POST("/api/v1/auth").
		WithJSON(map[string]string{
			"email":    "test@test.com",
			"password": "password",
			"action":   "login",
		}).
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("token")
}

func TestUserAPI(t *testing.T) {
	const USERID = 1

	e := httpexpect.New(t, PORT)

	// Test PATCH /api/v1/user
	t.Run("Update a user", func(t *testing.T) {
		updatedUser := map[string]interface{}{
			"email":           "test@test.com",
			"name":            "Test User",
			"userImageUri":    "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
			"companyName":     "Google",
			"companyImageUri": "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png",
		}

		e.PATCH("/api/v1/user").
			WithHeader("Authorization", "Bearer "+TOKEN).
			WithJSON(updatedUser).
			Expect().
			Status(200).
			JSON().Object().ContainsMap(updatedUser)
	})

	// Test GET /api/v1/user
	t.Run("Get all user", func(t *testing.T) {
		e.GET("/api/v1/user").
			WithHeader("Authorization", "Bearer "+TOKEN).
			Expect().
			Status(200).
			JSON().Array().NotEmpty()
	})
}

func TestDepartmentAPI(t *testing.T) {
	const DEPARTMENT_ID = 1
	e := httpexpect.New(t, PORT)

	// Test GET /api/v1/department
	t.Run("Get all departments", func(t *testing.T) {
		e.GET("/api/v1/department").
			WithHeader("Authorization", "Bearer "+TOKEN).
			Expect().
			Status(200).
			JSON().Array().Empty()
	})

	// Test POST /api/v1/department
	t.Run("Create a new department", func(t *testing.T) {
		department := map[string]interface{}{
			"name": "IT Department",
		}

		e.POST("/api/v1/department").
			WithHeader("Authorization", "Bearer "+TOKEN).
			WithJSON(department).
			Expect().
			Status(201).
			JSON().Object().
			ContainsMap(map[string]interface{}{
				"name": "IT Department",
			}).
			ContainsKey("departmentId").
			Value("departmentId").String().NotEmpty()
	})

	t.Run("Get all departments", func(t *testing.T) {
		e.GET("/api/v1/department").
			WithHeader("Authorization", "Bearer "+TOKEN).
			Expect().
			Status(200).
			JSON().Array().NotEmpty()
	})

	// Test PUT /api/v1/department/{id}
	t.Run("Update a department", func(t *testing.T) {
		departmentID := DEPARTMENT_ID
		updatedDepartment := map[string]interface{}{
			"name": "Updated IT Department",
		}

		e.PUT("/api/v1/department/{id}", departmentID).
			WithHeader("Authorization", "Bearer "+TOKEN).
			WithJSON(updatedDepartment).
			Expect().
			Status(200).
			JSON().Object().ContainsMap(updatedDepartment)
	})

	// Test DELETE /api/v1/department/{id}
	t.Run("Delete a department", func(t *testing.T) {
		departmentID := DEPARTMENT_ID

		e.DELETE("/api/v1/department/{id}", departmentID).
			WithHeader("Authorization", "Bearer "+TOKEN).
			Expect().
			Status(200)
	})
}

func TestEmployeeAPI(t *testing.T) {
	const EMPLOYEE_ID = "XX12345"

	e := httpexpect.New(t, PORT)

	getOrCreateDepartmentID := func() string {
		// Step 1: Try to get the list of departments
		departments := e.GET("/api/v1/department").
			WithHeader("Authorization", "Bearer "+TOKEN).
			Expect().
			Status(200).
			JSON().Array()

		// Step 2: If departments exist, return the first department's ID
		if departments.Length().Raw() > 0 {
			firstDepartment := departments.First().Object()
			departmentID := firstDepartment.Value("departmentId").Raw()

			// Handle the case where departmentId is a number or a string
			switch v := departmentID.(type) {
			case float64: // JSON numbers are unmarshalled as float64 by default
				return fmt.Sprintf("%d", int(v)) // Convert to string
			case int:
				return fmt.Sprintf("%d", v) // Convert to string
			case string:
				return v // Already a string
			default:
				panic(fmt.Sprintf("unexpected type for departmentId: %T", v))
			}
		}

		// Step 3: If no departments exist, create a new department
		newDepartment := map[string]interface{}{
			"name": "Default Department",
		}

		// Create the new department and get the response
		createdDepartment := e.POST("/api/v1/department").
			WithHeader("Authorization", "Bearer "+TOKEN).
			WithJSON(newDepartment).
			Expect().
			Status(201).
			JSON().Object()

		// Step 4: Handle the case where departmentId is a number or a string
		departmentID := createdDepartment.Value("departmentId").Raw()
		switch v := departmentID.(type) {
		case float64: // JSON numbers are unmarshalled as float64 by default
			return fmt.Sprintf("%d", int(v)) // Convert to string
		case int:
			return fmt.Sprintf("%d", v) // Convert to string
		case string:
			return v // Already a string
		default:
			panic(fmt.Sprintf("unexpected type for departmentId: %T", v))
		}
	}

	// Get or create a department ID
	DEPARTMENT_ID := getOrCreateDepartmentID()

	// Test POST /api/v1/emplyee
	t.Run("Create a new employee", func(t *testing.T) {
		employee := map[string]interface{}{
			"identityNumber":   EMPLOYEE_ID,
			"name":             "Bob Smith",
			"employeeImageUri": "1234",
			"gender":           "male",
			"departmentId":     DEPARTMENT_ID,
		}

		e.POST("/api/v1/employee").
			WithHeader("Authorization", "Bearer "+TOKEN).
			WithJSON(employee).
			Expect().
			Status(201).
			JSON().Object().
			ContainsMap(employee)
	})

	t.Run("Get all employees", func(t *testing.T) {
		e.GET("/api/v1/employee").
			WithHeader("Authorization", "Bearer "+TOKEN).
			Expect().
			Status(200).
			JSON().Array().NotEmpty()
	})

	// Test PATCH /api/v1/employee/{id}
	t.Run("Update a employee", func(t *testing.T) {
		updatedEmployee := map[string]interface{}{
			"identityNumber":   EMPLOYEE_ID,
			"name":             "Updated Bob Smith",
			"employeeImageUri": "1234",
			"gender":           "male",
			"departmentId":     DEPARTMENT_ID,
		}

		e.PATCH("/api/v1/employee/{id}", EMPLOYEE_ID).
			WithHeader("Authorization", "Bearer "+TOKEN).
			WithJSON(updatedEmployee).
			Expect().
			Status(200).
			JSON().Object().ContainsMap(updatedEmployee)
	})

	// Test DELETE /api/v1/employee/{id}
	t.Run("Delete a employee", func(t *testing.T) {
		e.DELETE("/api/v1/employee/{id}", EMPLOYEE_ID).
			WithHeader("Authorization", "Bearer "+TOKEN).
			Expect().
			Status(200)
	})
}
