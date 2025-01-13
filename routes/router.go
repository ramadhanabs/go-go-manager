package routes

import (
	"context"
	"database/sql"
	"go-go-manager/config"
	v1 "go-go-manager/controllers/v1"
	"go-go-manager/utils"
	"log"
	"net/http"
	"os"

	awsSdkCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func SetupRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
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
			// AWS credentials
			accessKey := cfg.AwsAccessKeyId
			secretKey := cfg.AwsSecretAccessKey
			region := cfg.S3Region
			bucketName := cfg.S3Bucket

			// Create custom credentials provider
			credProvider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

			// Load AWS configuration
			cfg, err := awsSdkCfg.LoadDefaultConfig(context.TODO(),
				awsSdkCfg.WithRegion(region),
				awsSdkCfg.WithCredentialsProvider(credProvider),
			)
			if err != nil {
				log.Fatalf("Unable to load SDK config: %v", err)
			}

			// Create S3 client
			client := s3.NewFromConfig(cfg)

			_, fileHeader, err := c.Request.FormFile("file")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read the file"})
				return
			}

			// Open the file
			file, err := os.Open(fileHeader.Filename)
			if err != nil {
				log.Fatalf("Unable to open file: %v", err)
			}
			defer file.Close()

			// Upload the file
			_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
				Bucket: &bucketName,
				Key:    &fileHeader.Filename,
				Body:   file,
			})
			if err != nil {
				log.Fatalf("Unable to upload file to S3: %v", err)
			}
		})
	}

	return router
}
