package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`               // Validates email format
	Password string `json:"password" binding:"required"`                  // Validates presence
	Action   string `json:"action" binding:"required,oneof=login signup"` // Validates specific values
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get users from v1"})
}

func AuthHandler(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.Action {
	case "login":
		// handle login here
	case "signup":
		// handle signup

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
	}
}
