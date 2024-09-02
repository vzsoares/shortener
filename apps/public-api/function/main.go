package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"apps/public-api/tools"

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

var apiUrl string
var remoteApi = "https://api-dev.zenhalab.com/shortener/v1"
var localApi = "http://localhost:4000"

func init() {
	if tools.DEBUG {
		apiUrl = localApi
	} else {
		apiUrl = remoteApi
	}

	http.HandleFunc(buildPath("/ping", nil), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(time.UnixDate)
	})

	http.HandleFunc(buildPath("/url/{id}", &GET), func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc(buildPath("/url", &POST), func(w http.ResponseWriter, r *http.Request) {})

	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpLambda.StripBasePath("/shortener/v1")
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
