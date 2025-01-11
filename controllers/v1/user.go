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
	print(user.Email)

	res := models.UserRequest{
		Email:           user.Email,
		Name:            user.Name,
		UserImageUri:    user.UserImageUri,
		CompanyName:     user.CompanyName,
		CompanyImageUri: user.CompanyImageUri,
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

	models.UpdateProfile(body, v.UserID)

	c.JSON(200, gin.H{"message": "Profile updated successfully"})
}
