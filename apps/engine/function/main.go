package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"apps/engine/domain"
	"apps/engine/handler"
	"apps/engine/store"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var GET = "GET"
var PUT = "PUT"
var POST = "POST"
var DELETE = "DELETE"

var httpLambda *httpadapter.HandlerAdapter

func buildPath(p string, m *string) string {
	var res string
	basePath := "/engine"
	if m == nil {
		res = fmt.Sprint(basePath, p)
	} else {
		res = fmt.Sprint(*m, " ", basePath, p)
	}
	return res
}

var apiUrl = "https://dynamodb.us-east-1.amazonaws.com"
var apiUrlLocal = "http://host.docker.internal:8000"

func init() {
	http.HandleFunc(buildPath("/ping", nil), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(time.UnixDate)
	})
	ctx := context.TODO()
	store := store.NewDynamoStore(ctx, apiUrlLocal, true)
	domain := domain.NewUrlDomain(ctx, store)
	handler := handler.NewHttpHandler(ctx, domain)

	http.HandleFunc(buildPath("/url/{id}", &GET), handler.GetHandler)
	http.HandleFunc(buildPath("/url/{id}", &DELETE), handler.DeleteHandler)
	http.HandleFunc(buildPath("/url", &POST), handler.PostHandler)

	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
