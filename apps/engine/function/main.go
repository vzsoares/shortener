package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"apps/engine/domain"
	"apps/engine/handlers"
	"apps/engine/store"
	"apps/engine/tools"

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
var skipHttps bool

func init() {
	if tools.DEBUG {
		apiUrl = tools.LOCAL_DB_URL_DOCKER
		skipHttps = true
	} else {
		apiUrl = tools.REMOTE_DB_URL
		skipHttps = false
	}

	http.HandleFunc(buildPath("/ping", nil), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(time.UnixDate)
	})
	ctx := context.TODO()
	store := store.NewDynamoStore(ctx, apiUrl, skipHttps)
	domain := domain.NewUrlDomain(ctx, store)
	handler := handlers.NewHttpHandler(ctx, domain)

	http.HandleFunc(buildPath("/url/{id}", &GET),
		handlers.AuthMiddleware(handler.GetHandler))
	http.HandleFunc(buildPath("/url/{id}", &DELETE),
		handlers.AuthMiddleware(handler.DeleteHandler))
	http.HandleFunc(buildPath("/url", &POST),
		handlers.AuthMiddleware(handler.PostHandler))

	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpLambda.StripBasePath("/shortener/v1")
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
