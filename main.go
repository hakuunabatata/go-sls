package main

import (
	"encoding/json"
	"strconv"

	events "github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	aws "github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	uuid "github.com/google/uuid"
)

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func InsertProduct(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var product Product
	err := json.Unmarshal([]byte(request.Body), &product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	product.ID = uuid.New().String()

	session := session.Must(session.NewSession())

	service := dynamodb.New(session)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("ProductsVideo"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(product.ID),
			},
			"name": {
				S: aws.String(product.Name),
			},
			"price": {
				N: aws.String(strconv.Itoa(product.Price)),
			},
		},
	}

	_, err = service.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := json.Marshal(product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func main() {
	lambda.Start(InsertProduct)
}
