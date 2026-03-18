package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucketName string

func init() {
	log.Println("=== Initializing with AWS SDK v2 ===")
	
	bucketName = os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("BUCKET_NAME environment variable not set")
	}
	log.Printf("BUCKET_NAME = %s", bucketName)
	
	// Load AWS config with explicit region
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-1"),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config: %v", err)
	}
	
	// Create S3 client
	s3Client = s3.NewFromConfig(cfg)
	log.Println("S3 client created successfully")
}

func handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("=== Handler Invoked ===")
	
	key := request.QueryStringParameters["key"]
	log.Printf("Requested key: %s", key)

	if key == "" {
		log.Println("Missing key parameter")
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "missing key",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	log.Printf("Getting object s3://%s/%s", bucketName, key)
	
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	
	if err != nil {
		log.Printf("Error getting object: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error: %v", err),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error reading body: %v", err),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}
	
	log.Printf("Success! Read %d bytes", len(body))
	
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}