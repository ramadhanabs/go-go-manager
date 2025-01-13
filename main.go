package main

import (
	"context"
	"fmt"
	"go-go-manager/config"
	"go-go-manager/db"
	"go-go-manager/routes"
	"log"

	awsSdkCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	cfg := config.LoadConfig()

	db.InitDB(cfg)
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
		log.Println("Database connection closed.")
	}()

	// AWS credentials
	accessKey := cfg.AwsAccessKeyId
	secretKey := cfg.AwsSecretAccessKey
	region := cfg.S3Region
	bucketName := cfg.S3Bucket

	// Create custom credentials provider
	credProvider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	// Load AWS configuration
	awsCfg, err := awsSdkCfg.LoadDefaultConfig(context.TODO(),
		awsSdkCfg.WithRegion(region),
		awsSdkCfg.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(awsCfg)

	r := routes.SetupRouter(cfg, db.DB, s3Client, bucketName)

	fmt.Printf("Starting server on port %s...\n", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
