package main

import (
	"bytes"
	"encoding/json"
	_ "github.com/aws/aws-xray-sdk-go/xray"
	"log"
	"strings"
	"time"

	"github.com/localz/go-lambda-example/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Println(req.PathParameters)
	projectID, exists := req.PathParameters["projectID"]

	if !exists {
		return events.APIGatewayProxyResponse{Body: "projectID required", StatusCode: 400}, nil
	}

	resourceID, exists := req.PathParameters["resourceID"]

	if !exists {
		return events.APIGatewayProxyResponse{Body: "resourceID required", StatusCode: 400}, nil
	}

	filenames, exists := req.QueryStringParameters["filenames"]

	if !exists {
		return events.APIGatewayProxyResponse{Body: "filenames required", StatusCode: 400}, nil
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	if err != nil {
		log.Println(err)
	}

	svc := s3.New(sess)

	var urls []string

	for _, filename := range strings.Split(filenames, ",") {
		s3req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(repository.GetBucket()),
			Key:    aws.String(projectID + "/" + resourceID + "/" + filename),
		})

		str, err := s3req.Presign(15 * time.Minute)

		if err != nil {
			log.Println(err)
		}

		urls = append(urls, str)
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
