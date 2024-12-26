package files

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:embed static/* templates/*
var Files embed.FS

//go:embed migrations/*
var Migrations embed.FS

var s3client *s3.Client

func InitS3Client() error {
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "eu-central-1"
	}
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint == "" {
		awsEndpoint = "http://localhost:4566"
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return err
	}

	s3client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(awsEndpoint)
	})

	return nil
}

func UploadProfilePic(profilePicName string, profilePicData []byte) error {
	_, err := s3client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("profile-pictures"),
		Key:    aws.String(profilePicName),
		Body:   bytes.NewReader(profilePicData),
	})
	return err
}

func DeleteProfilePic(profilePicName string) error {
	_, err := s3client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String("profile-pictures"),
		Key:    aws.String(profilePicName),
	})
	return err
}

func GetProfilePicPath(profilePicName sql.NullString) string {
	if profilePicName.Valid {
		return "https://localhost.localstack.cloud:4566/profile-pictures/" + profilePicName.String
	}
	return "/static/img/blank-profile-picture.png"
}
