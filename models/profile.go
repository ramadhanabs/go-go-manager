package models

import (
	"database/sql"
	"go-go-manager/db"
)

type UserRequest struct {
	Email           string `json:"email" binding:"required,email"`
	Name            string `json:"name" binding:"required,min=4,max=52"`
	UserImageUri    string `json:"userImageUri" binding:"omitempty,uri"`
	CompanyName     string `json:"companyName" binding:"required,min=4,max=52"`
	CompanyImageUri string `json:"companyImageUri" binding:"omitempty,uri"`
}

func CheckEmailDuplicate(email string, userID uint) (bool, error) {
	query := "SELECT id FROM users WHERE email = $1 AND id != $2"
	var id uint
	err := db.DB.QueryRow(query, email, userID).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
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
