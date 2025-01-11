package v1

import (
	"context"
	"fmt"
	"go-go-manager/config"
	"go-go-manager/utils"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsSdkCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
)

const maxFileSize = 100 * 1024 // 100 KB

type FileHandler struct {
	s3Client *s3.Client
	s3Bucket string
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	awsCfg, err := awsSdkCfg.LoadDefaultConfig(
		context.Background(),
		awsSdkCfg.WithRegion(cfg.S3Region),
		awsSdkCfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AwsAccessKeyId,
				cfg.AwsSecretAccessKey,
				"",
			),
		),
	)
	if err != nil {
		log.Fatalf("cannot load the AWS configs: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// o.UsePathStyle = true
		// o.BaseEndpoint = aws.String(cfg.S3Endpoint)
	})

	return &FileHandler{
		s3Client: s3Client,
		s3Bucket: os.Getenv("AWS_S3_BUCKET_NAME"),
	}
}

type FileRequest struct {
	File string `json:"file" binding:"required"`
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	jwtClaims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read the file"})
		return
	}

	if !isValidFileType(fileHeader.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed types: jpeg, jpg, png"})
		return
	}

	if fileHeader.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the maximum limit of 100 KiB"})
		return
	}

	uri, err := h.uploadToS3(jwtClaims.Email, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file to S3"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uri": uri,
	})
}

func isValidFileType(filename string) bool {
	ext := strings.ToLower(strings.TrimSpace(filename[strings.LastIndex(filename, ".")+1:]))
	return ext == "jpeg" || ext == "jpg" || ext == "png"
}

func (h *FileHandler) uploadToS3(email string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String(h.s3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", email, fileHeader.Filename)),
		ACL:    types.ObjectCannedACLPublicRead,
		Body:   file,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute*2)
	defer cancel()

	_, err = h.s3Client.PutObject(ctx, input)
	if err != nil {
		return "", err
	}

	s3URI := fmt.Sprintf("s3://%s/%s/%s", h.s3Bucket, email, fileHeader.Filename)
	return s3URI, nil
}
