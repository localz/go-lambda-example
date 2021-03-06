package main

import (
	"encoding/json"

  _ "github.com/aws/aws-xray-sdk-go/xray"

  // Importing the plugins enables collection of AWS resource information at runtime.
  // Every plugin should be imported after "github.com/aws/aws-xray-sdk-go/xray" library.

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/localz/go-lambda-example/repository"
)

// PersonResponse
type PersonResponse struct {
	Person repository.Person `json:"data"`
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	person, err := repository.GetPerson(req.PathParameters["id"])

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "User not found", StatusCode: 404}, nil
	}

	body, err := json.Marshal(PersonResponse{*person})

	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Unable to marshal JSON", StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
