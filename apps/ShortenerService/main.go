package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

type Response struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

var httpLambda *httpadapter.HandlerAdapter

func init() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(Response{From: "net/http", Message: time.Now().Format(time.UnixDate)+"123456"})
	})

	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func main() {
	lambda.Start(Handler)

}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return httpLambda.ProxyWithContext(ctx, req)
}
