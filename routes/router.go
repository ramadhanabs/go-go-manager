package routes

import (
	"context"
	"database/sql"
	"go-go-manager/config"
	v1 "go-go-manager/controllers/v1"
	"go-go-manager/utils"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupRouter(cfg *config.Config, db *sql.DB, s3Client *s3.Client, bucketName string) *gin.Engine {
	router := gin.Default()

	employeeHandler := v1.NewEmployeeHandler(db)
	// v1FileHandler := v1.NewFileHandler(cfg)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("isImage", utils.IsImageURI)
	}

	v1Group := router.Group("/api/v1")
	{
		v1Group.POST("/auth", v1.AuthHandler)
		v1Group.GET("/user", v1.GetUsers)
		v1Group.PATCH("/user", v1.UpdateUser)
		v1Group.POST("/department", v1.CreateDepartment)
		v1Group.GET("/department", v1.GetDepartments)
		v1Group.PATCH("/department/:departmentId", v1.UpdateDepartment)
		v1Group.DELETE("/department/:departmentId", v1.DeleteDepartment)

		// Employee routes
		v1Group.POST("/employee", employeeHandler.CreateEmployee())
		v1Group.GET("/employee", employeeHandler.GetEmployees())
		v1Group.PATCH("/employee/:identityNumber", employeeHandler.UpdateEmployee())
		v1Group.DELETE("/employee/:identityNumber", employeeHandler.DeleteEmployee())

		// v1Group.POST("/file", v1FileHandler.UploadFile)
		v1Group.POST("/file", func(c *gin.Context) {
			_, fileHeader, err := c.Request.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read the file"})
				return
			}

			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				log.Fatalf("Unable to open file: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			defer file.Close()

			// Upload the file
			_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucketName,
				Key:    &fileHeader.Filename,
				Body:   file,
			})
			if err != nil {
				log.Fatalf("Unable to upload file to S3: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		})
	}

	return router
}
