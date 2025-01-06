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
