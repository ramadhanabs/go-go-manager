package models

import (
	"database/sql"
	"fmt"
	"go-go-manager/db"
)

type Department struct {
	ID        uint
	Name      string
	CreatedAt string
	UpdatedAt string
}

func CreateDepartment(name string) (Department, error) {
	query := "INSERT INTO department (name) VALUES ($1) RETURNING id, name"

	var department Department
	err := db.DB.QueryRow(query, name).Scan(&department.ID, &department.Name)
	if err != nil {
		return Department{}, fmt.Errorf("failed to create department: %v", err)
	}

	return department, nil
}

func FindDepartmentByName(name string) (Department, error) {
	query := "SELECT id, name FROM department WHERE name = $1"
	var department Department

	row := db.DB.QueryRow(query, name)

	err := row.Scan(&department.ID, &department.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Department not found")
			return Department{}, fmt.Errorf("no department found with name %s", name)
		}
		return Department{}, err
	}

	return department, nil
}

func UpdateDepartment(id string, name string) (Department, error) {
	query := "UPDATE department SET name = $1 WHERE id = $2 RETURNING id, name"
	var department Department
	err := db.DB.QueryRow(query, name, id).Scan(&department.ID, &department.Name)
	if err != nil {
		return Department{}, fmt.Errorf("failed to update department: %v", err)
	}

	return department, err
}

func FindDepartmentById(userID uint, id string) (Department, error) {
	query := "SELECT id, name FROM department WHERE id = $1 AND userId = $2"
	var department Department

	row := db.DB.QueryRow(query, id, userID)

	err := row.Scan(&department.ID, &department.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Department not found")
			return Department{}, fmt.Errorf("no department found with id %s", id)
		}
		return Department{}, err
	}

	return department, nil
}

func DeleteDepartment(id string) error {
	query := "DELETE FROM department WHERE id = $1"

	result, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department with id %s not found", id)
	}

	return nil
}

func CountEmployeesByDepartment(departmentId string) (int, error) {
	query := "SELECT COUNT(*) FROM employee WHERE department_id = $1"
	var count int
	err := db.DB.QueryRow(query, departmentId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
