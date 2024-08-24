package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableName = "sample"
var BaseEndpoint = "http://localhost:8000/"

func Handler(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointOptions.DisableHTTPS = true
		o.BaseEndpoint = aws.String(BaseEndpoint)
	})

	input := &dynamodb.GetItemInput{
		TableName: &TableName,
		Key: map[string]ddbtypes.AttributeValue{
			"city": &ddbtypes.AttributeValueMemberS{Value: "Mumbai"},
			"Name": &ddbtypes.AttributeValueMemberS{Value: "Ali"},
		}}

	res, err := client.GetItem(ctx, input)
	if err != nil {
		panic(err.Error())
	}

	body, err := json.Marshal(res)
	if err != nil {
		panic(err.Error())
	}

	qinput := &dynamodb.QueryInput{
		TableName:              &TableName,
		KeyConditionExpression: aws.String("#name = :name"),
		ExpressionAttributeValues: map[string]ddbtypes.AttributeValue{
			":name": &ddbtypes.AttributeValueMemberS{Value: "Ali"},
		},
		ExpressionAttributeNames: map[string]string{
			"#name": "Name",
		},
	}

	qres, err := client.Query(ctx, qinput)
	if err != nil {
		panic(err.Error())
	}

	qbody, err := json.Marshal(qres)
	if err != nil {
		panic(err.Error())
	}
	println(string(qbody))
	println()

	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: string(body)}, nil
}
