package main

import (
	"context"
	"encoding/json"
	"fmt"
	"libs/utils"
	"log"
	"net/http"
	"time"

	"apps/public-api/handlers"
	"apps/public-api/tools"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var GET = "GET"
var PUT = "PUT"
var POST = "POST"
var DELETE = "DELETE"
var OPTIONS = "OPTIONS"

var httpLambda *httpadapter.HandlerAdapter

func buildPath(p string, m *string) string {
	var res string
	basePath := "/public-api"
	if m == nil {
		res = fmt.Sprint(basePath, p)
	} else {
		res = fmt.Sprint(*m, " ", basePath, p)
	}
	return res
}

var apiUrl string
var apiUrlRemote string = "https://api-dev.zenhalab.com/shortener/v1"
var apiUrlLocal string = "http://localhost:4000"
var apiKeyA4 string
var parameterStore *utils.Ssm

func init() {
	if tools.DEBUG {
		apiUrl = apiUrlLocal
	} else {
		apiUrl = apiUrlRemote
	}

	client := http.Client{}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}
	if parameterStore == nil {
		parameterStore = utils.NewSmmStore(cfg, ctx)

	}

	if apiKeyA4 == "" {
		apiKeyA4 = parameterStore.Get("API_KEY_A4")

	}

	handler := handlers.NewHttpHandler(ctx, client, apiUrl, apiKeyA4)

	http.HandleFunc(buildPath("/ping", &GET), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(time.UnixDate)
	})

	http.HandleFunc(buildPath("/url/{id}", &GET), handler.GetHandler)
	http.HandleFunc(buildPath("/url", &POST), handler.PostHandler)

	http.HandleFunc(buildPath("/", &OPTIONS), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	})

	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpLambda.StripBasePath("/shortener/v1")
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	if tools.DEBUG {
		println("Running debug server...")
		log.Fatal(http.ListenAndServe(":3000", nil))
	} else {
		lambda.Start(Handler)
	}
}
