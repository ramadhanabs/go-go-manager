package models

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

type Employee struct {
	IdentityNumber   string `json:"identityNumber" binding:"required,min=5,max=33"`
	Name             string `json:"name" binding:"required,min=4,max=33"`
	Gender           Gender `json:"gender" binding:"required"` // Enum: "male" or "female"
	DepartmentID     string `json:"departmentId" binding:"required"`
	EmployeeImageURI string `json:"employeeImageUri" binding:"required,uri,isImage"` // New field
}
