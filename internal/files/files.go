package files

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

	buckets := []string{"profile-pictures", "vajb-pictures"}
	for _, bucket := range buckets {
		exists, err := checkBucketExists(s3client, bucket)
		if err != nil {
			return err
		}
		if !exists {
			err = createBucket(bucket, awsRegion)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func checkBucketExists(client *s3.Client, bucketName string) (bool, error) {
	_, err := client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func createBucket(bucketName, region string) error {
	_, err := s3client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	return err
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

func UploadVajbPic(vajbPicName string, vajbPicData []byte) error {
	_, err := s3client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("vajb-pictures"),
		Key:    aws.String(vajbPicName),
		Body:   bytes.NewReader(vajbPicData),
	})
	return err
}
