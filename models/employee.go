package models

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Employee struct {
	IdentityNumber   string `json:"identityNumber"`
	Name             string `json:"name"`
	Gender           Gender `json:"gender"` // Enum: "male" or "female"
	DepartmentID     string `json:"departmentId"`
	EmployeeImageURI string `json:"employeeImageUri"` // New field
}