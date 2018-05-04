package main

import (
	"bytes"
	"encoding/json"
	_ "github.com/aws/aws-xray-sdk-go/xray"
	"log"
	"strings"

	"github.com/localz/go-lambda-example/repository"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	var results []string

	for _, filename := range strings.Split(filenames, ",") {
		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(repository.GetBucket()),
			Key:    aws.String(projectID + "/" + resourceID + "/" + filename),
		})

		if err != nil {
			log.Println(err)
			results = append(results, err.Error())
		}
	}

	if len(results) == 0 {
		return events.APIGatewayProxyResponse{Body: "[]", StatusCode: 200}, nil
	}

	buffer := &bytes.Buffer{}
	e := json.NewEncoder(buffer)
	e.SetEscapeHTML(false)
	e.Encode(results)

	return events.APIGatewayProxyResponse{Body: buffer.String(), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
