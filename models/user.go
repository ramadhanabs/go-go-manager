package models

import (
	"database/sql"
	"fmt"
	"go-go-manager/db"
)

type User struct {
	ID              uint
	Email           string
	Name            string
	Password        string
	UserImageUri    string
	CompanyName     string
	CompanyImageUri string
	CreatedAt       string
	UpdatedAt       string
}

func FindUserByEmail(email string) (User, error) {
	query := "SELECT id, email, password FROM users WHERE email = $1"
	var user User

	row := db.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("no user found with email: %s", email)
		}
		return User{}, err
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	query := "SELECT id, email, name, user_image_uri, company_name, company_image_uri FROM users WHERE id = $1"
	var user User

	row := db.DB.QueryRow(query, id)

	err := row.Scan(&user.ID, &user.Email, &user.Name, &user.UserImageUri, &user.CompanyName, &user.CompanyImageUri)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("no user found with id: %d", id)
		}
		return User{}, err
	}

	return user, nil
}
