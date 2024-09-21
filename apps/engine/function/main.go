package main

import (
	"context"
	"encoding/json"
	"fmt"
	"libs/utils"
	"log"
	"net/http"
	"time"

	"apps/engine/domain"
	"apps/engine/handlers"
	"apps/engine/store"
	"apps/engine/tools"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
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
var parameterStore *utils.Ssm

func init() {
	if tools.DEBUG {
		apiUrl = tools.LOCAL_DB_URL
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
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}

	store := store.NewDynamoStore(ctx, apiUrl, skipHttps, cfg)
	domain := domain.NewUrlDomain(ctx, store)
	handler := handlers.NewHttpHandler(ctx, domain)
	if parameterStore == nil {
		parameterStore = utils.NewSmmStore(cfg, ctx)
	}

	http.HandleFunc(buildPath("/url/{id}", &GET),
		handlers.AuthMiddleware(handler.GetHandler, parameterStore))
	http.HandleFunc(buildPath("/url/{id}", &DELETE),
		handlers.AuthMiddleware(handler.DeleteHandler, parameterStore))
	http.HandleFunc(buildPath("/url", &POST),
		handlers.AuthMiddleware(handler.PostHandler, parameterStore))

	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpLambda.StripBasePath("/shortener/v1")
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	if tools.DEBUG {
		println("Running debug server...")
		log.Fatal(http.ListenAndServe(":4000", nil))
	} else {
		lambda.Start(Handler)
	}
}
