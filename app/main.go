
package main

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucket string

func init() {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	s3Client = s3.NewFromConfig(cfg)

	bucket = os.Getenv("BUCKET_NAME")
}

func handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	key := request.QueryStringParameters["key"]

	obj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 404,
			Body:       "Object not found",
		}, nil
	}

	defer obj.Body.Close()

	data, _ := io.ReadAll(obj.Body)

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type": "application/octet-stream",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
