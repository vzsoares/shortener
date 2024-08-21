package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	lambda.Start(handler)
}

const tableName = "ursldas"

func handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	cfg, _ := config.LoadDefaultConfig(ctx)

	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.GetItemInput{}

	client.GetItem(ctx, input)

	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: ""}, nil
}
