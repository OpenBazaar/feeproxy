package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/OpenBazaar/feeproxy"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	lambda.Start(updateS3)
}

func updateS3() error {
	feeData, err := feeproxy.Query()
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = writeToS3(
		os.Getenv("AWS_S3_REGION"),
		os.Getenv("AWS_S3_BUCKET"),
		os.Getenv("AWS_S3_FILENAME"),
		feeData,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func writeToS3(region, bucket, filename string, data []byte) error {
	if filename == "" {
		filename = "fees"
	}

	_, err := s3.New(
		session.New(),
		aws.NewConfig().WithRegion(region),
	).PutObject(&s3.PutObjectInput{
		Key:           aws.String(filename),
		Bucket:        aws.String(bucket),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String("application/json"),
	})

	return err
}
