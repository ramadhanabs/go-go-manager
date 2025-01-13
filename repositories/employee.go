package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-go-manager/models"
	"strconv"
)

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) AddEmployee(employee models.Employee) error {
	query := `
		INSERT INTO employees (identity_number, name, gender, department_id, employee_image_uri)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.DB.ExecContext(context.Background(), query,
		employee.IdentityNumber,
		employee.Name,
		employee.Gender,
		employee.DepartmentID,
		employee.EmployeeImageURI,
	)
	return err
}

func (r *EmployeeRepository) GetEmployeeByIdentityNumber(identityNumber string) (*models.Employee, error) {
	query := `
		SELECT identity_number, name, gender, department_id, employee_image_uri
		FROM employees
		WHERE identity_number = $1
	`
	var employee models.Employee
	err := r.DB.QueryRowContext(context.Background(), query, identityNumber).Scan(
		&employee.IdentityNumber,
		&employee.Name,
		&employee.Gender,
		&employee.DepartmentID,
		&employee.EmployeeImageURI,
	)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) UpdateEmployee(identityNumber string, updatedEmployee models.Employee) error {
	query := `
		UPDATE employees
		SET name = $1, gender = $2, department_id = $3, employee_image_uri = $4
		WHERE identity_number = $5
	`
	_, err := r.DB.ExecContext(context.Background(), query,
		updatedEmployee.Name,
		updatedEmployee.Gender,
		updatedEmployee.DepartmentID,
		updatedEmployee.EmployeeImageURI,
		identityNumber,
	)
	return err
}

func (r *EmployeeRepository) DeleteEmployee(identityNumber string) error {
	query := `
		DELETE FROM employees
		WHERE identity_number = $1
	`
	_, err := r.DB.ExecContext(context.Background(), query, identityNumber)
	return err
}

func (r *EmployeeRepository) FilterEmployees(filters map[string]string) ([]models.Employee, error) {
	query := `
		SELECT identity_number, name, gender, department_id, employee_image_uri
		FROM employees
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if identityNumber, ok := filters["identityNumber"]; ok {
		query += fmt.Sprintf(" AND identity_number LIKE $%d || '%%'", argCount)
		args = append(args, identityNumber)
		argCount++
	}
	if name, ok := filters["name"]; ok {
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+name+"%")
		argCount++
	}
	if gender, ok := filters["gender"]; ok {
		query += fmt.Sprintf(" AND gender = $%d", argCount)
		args = append(args, gender)
		argCount++
	}
	if departmentID, ok := filters["departmentId"]; ok {
		query += fmt.Sprintf(" AND department_id = $%d", argCount)
		args = append(args, departmentID)
		argCount++
	}

	// Add LIMIT and OFFSET for pagination
	limit, _ := strconv.Atoi(filters["limit"])
	offset, _ := strconv.Atoi(filters["offset"])
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, limit)
		argCount++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, offset)
		argCount++
	}

	rows, err := r.DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.IdentityNumber,
			&emp.Name,
			&emp.Gender,
			&emp.DepartmentID,
			&emp.EmployeeImageURI,
		)
		if err != nil {
			return nil, err
		}
		employees = append(employees, emp)
	}

	return employees, nil
}
