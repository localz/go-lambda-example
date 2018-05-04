package main

import (
	"bytes"
	"encoding/json"
	_ "github.com/aws/aws-xray-sdk-go/xray"
	"log"
	"time"

	"github.com/localz/go-lambda-example/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projectId, exists := req.PathParameters["projectId"]

	if !exists {
		return events.APIGatewayProxyResponse{Body: "projectId required", StatusCode: 400}, nil
	}

	resourceId, exists := req.PathParameters["resourceId"]

	if !exists {
		return events.APIGatewayProxyResponse{Body: "resourceId required", StatusCode: 400}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})

	svc := s3.New(sess)

	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(repository.GetBucket()),
		Prefix: aws.String(projectId + "/" + resourceId),
	})
	if err != nil {
		log.Println(err)
	}

	var urls []string

	for _, item := range resp.Contents {
		if (*item.Key)[len(*item.Key)-1:] == "/" {
			continue
		}

		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(repository.GetBucket()),
			Key:    item.Key,
		})

		urlStr, err := req.Presign(15 * time.Minute)

		if err != nil {
			log.Println("Failed to sign request", err)
		}

		urls = append(urls, urlStr)
	}

	if len(urls) == 0 {
		return events.APIGatewayProxyResponse{Body: "[]", StatusCode: 200}, nil
	}

	buffer := &bytes.Buffer{}
	e := json.NewEncoder(buffer)
	e.SetEscapeHTML(false)
	e.Encode(urls)

	return events.APIGatewayProxyResponse{Body: buffer.String(), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
