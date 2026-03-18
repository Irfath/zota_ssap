package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3Client *s3.S3
var bucketName string

func init() {
	log.Println("=== Initializing with AWS SDK v1 ===")
	
	bucketName = os.Getenv("BUCKET_NAME")
	log.Printf("BUCKET_NAME = %s", bucketName)
	
	// Create session with explicit region
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
	}))
	
	s3Client = s3.New(sess)
	log.Println("S3 client created")
}

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("=== Handler Invoked ===")
	
	key := request.QueryStringParameters["key"]
	log.Printf("Key: %s", key)

	if key == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "missing key",
		}, nil
	}

	log.Printf("Getting object: %s/%s", bucketName, key)
	
	resp, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	
	if err != nil {
		log.Printf("Error: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error: %v", err),
		}, nil
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read error: %v", err)
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Read error: %v", err),
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