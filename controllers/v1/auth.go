package v1

import (
	"go-go-manager/models"
	"go-go-manager/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`               // Validates email format
	Password string `json:"password" binding:"required,min=8,max=32"`     // Validates presence
	Action   string `json:"action" binding:"required,oneof=create login"` // Validates specific values
}

func AuthHandler(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.Action {
	case "login":
		user, err := models.FindUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}

		// check password validity
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password mismatch"})
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"email": user.Email,
			"token": token,
		})
	case "signup":
		// handle signup
		_, err := models.FindUserByEmail(req.Email)
		if err != nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hashing password"})
				return
			}

			user, err := models.CreateUser(req.Email, string(hashedPassword))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			token, err := utils.GenerateJWT(user.ID, user.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"email": user.Email,
				"token": token,
			})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
	}
}
