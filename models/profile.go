package models

import (
	"go-go-manager/db"
)

type UserRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Name            string `json:"name" binding:"required"`
	UserImageUri    string `json:"userImageUri"`
	CompanyName     string `json:"companyName"`
	CompanyImageUri string `json:"companyImageUri"`
}

func CheckEmailDuplicate(email string) (bool, error) {

	query := "SELECT email FROM users WHERE email = $1"
	var userEmail string
	err := db.DB.QueryRow(query, email).Scan(&userEmail)

	if err != nil {
		return false, err
	}

	return true, nil

}

func UpdateProfile(req UserRequest, id uint) (UserRequest, error) {

	query := "UPDATE users SET email = $1, name = $2, user_image_uri = $3, company_name = $4, company_image_uri = $5 WHERE id = $6"
	_, err := db.DB.Exec(query, req.Email, req.Name, req.UserImageUri, req.CompanyName, req.CompanyImageUri, id)

	if err != nil {
		return UserRequest{}, err
	}
	return req, nil

}
