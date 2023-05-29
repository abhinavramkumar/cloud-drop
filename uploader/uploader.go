package uploader

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	Client         *s3.Client
	Bucket         string
	OsFilePath     string
	UploadFilePath string
}

func (s3Client *S3Client) CreateClient(profileName string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(profileName))
	if err != nil {
		log.Fatal(err)
	}

	s3Client.Client = s3.NewFromConfig(cfg)
}

func (s3Client *S3Client) UploadToS3() (bool, error) {
	file, err := os.Open(s3Client.OsFilePath)
	if err != nil {
		println(err)
		return false, err
	}
	defer file.Close()

	params := &s3.PutObjectInput{
		Bucket:       aws.String(s3Client.Bucket),
		Key:          aws.String(s3Client.UploadFilePath),
		Body:         file,
		ACL:          types.ObjectCannedACLPublicRead,
		CacheControl: aws.String("public, max-age=31536000"),
	}

	// Upload the file to S3
	_, err = s3Client.Client.PutObject(context.TODO(), params)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return false, err
	}

	return true, nil
}
