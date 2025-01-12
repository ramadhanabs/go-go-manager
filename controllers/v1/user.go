package v1

import (
	"go-go-manager/models"
	"go-go-manager/utils"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	auth = auth[7:]

	v, err := utils.ValidateJWT(auth)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	user, err := models.FindUserById(v.UserID)

	res := models.UserRequest{
		Email:           user.Email,
		Name:            user.Name.String,
		UserImageUri:    user.UserImageUri.String,
		CompanyName:     user.CompanyName.String,
		CompanyImageUri: user.CompanyImageUri.String,
	}

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, res)

}

func UpdateUser(c *gin.Context) {

	auth := c.GetHeader("Authorization")
	auth = auth[7:]

	v, err := utils.ValidateJWT(auth)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	var body models.UserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if body.Email == "" || body.Name == "" || body.CompanyName == "" || body.UserImageUri == "" || body.CompanyImageUri == "" {
		c.JSON(400, gin.H{"error": "All fields are required"})
		return
	}

	ed, _ := models.CheckEmailDuplicate(body.Email)

	if ed {
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	models.UpdateProfile(body, v.UserID)

	c.JSON(200, gin.H{"message": "Profile updated successfully"})
}
