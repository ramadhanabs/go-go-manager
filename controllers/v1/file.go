package v1

import (
	"context"
	"fmt"
	"go-go-manager/config"
	"go-go-manager/utils"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsSdkCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

const maxFileSize = 100 * 1024 // 100 KB

type FileHandler struct {
	s3Client *s3.Client
	s3Bucket string
	uploader *manager.Uploader
}

func NewFileHandler(cfg *config.Config) *FileHandler {
	awsCfg, err := awsSdkCfg.LoadDefaultConfig(
		context.TODO(),
		awsSdkCfg.WithRegion(cfg.S3Region),
		awsSdkCfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AwsAccessKeyId,
				cfg.AwsSecretAccessKey,
				"",
			),
		),
		awsSdkCfg.WithClientLogMode(aws.LogRequest|aws.LogResponse), // dev mode
	)
	if err != nil {
		log.Fatalf("cannot load the AWS configs: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		// For testing purpose with Localstack
		// o.UsePathStyle = true
		// o.BaseEndpoint = aws.String(cfg.S3Endpoint)
	})

	return &FileHandler{
		uploader: manager.NewUploader(s3Client),
		s3Client: s3Client,
		s3Bucket: cfg.S3Bucket,
	}
}

type FileRequest struct {
	File string `json:"file" binding:"required"`
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing"})
		return
	}

	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return
	}

	auth = auth[7:]
	v, err := utils.ValidateJWT(auth)
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

	uri, err := h.uploadToS3(v.Email, fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uri": uri,
	})
}

func (h *FileHandler) uploadToS3(email string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// input := &s3.PutObjectInput{
	// 	Bucket: aws.String(h.s3Bucket),
	// 	Key:    aws.String(fmt.Sprintf("%s/%s", email, fileHeader.Filename)),
	// 	ACL:    types.ObjectCannedACLPublicRead,
	// 	Body:   file,
	// }

	contentType := getContentType(fileHeader.Filename)

	_, err = h.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      &h.s3Bucket,
		Key:         aws.String(fmt.Sprintf("%s/%s", email, fileHeader.Filename)),
		Body:        file,
		ContentType: &contentType,
	})
	if err != nil {
		return "", err
	}

	s3URI := fmt.Sprintf("s3://%s/%s/%s", h.s3Bucket, email, fileHeader.Filename)
	return s3URI, nil
}

func isValidFileType(filename string) bool {
	ext := strings.ToLower(strings.TrimSpace(filename[strings.LastIndex(filename, ".")+1:]))
	return ext == "jpeg" || ext == "jpg" || ext == "png"
}

func getContentType(filename string) string {
	ext := filename[strings.LastIndex(filename, ".")+1:]
	contentTypeMap := map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
		"gif":  "image/gif",
		"pdf":  "application/pdf",
		"txt":  "text/plain",
	}

	contentType := contentTypeMap[ext]
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return contentType
}
