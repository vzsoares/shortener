package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	etools "apps/engine/tools"
	"apps/engine/types"
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

var apiUrl string = "https://api-dev.zenhalab.com/shortener/v1"
var apiKeyA4 string
var parameterStore *etools.Ssm

func init() {
	if tools.DEBUG {
	} else {
	}

	client := http.Client{}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}
	if parameterStore == nil {
		parameterStore = etools.NewSmmStore(cfg, ctx)

	}

	if apiKeyA4 == "" {
		apiKeyA4 = parameterStore.Get("API_KEY_A4")

	}

	http.HandleFunc(buildPath("/ping", nil), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(time.UnixDate)
	})

	http.HandleFunc(buildPath("/url/{id}", &GET), func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			// TODO error page
			panic("No id")
		}
		request, err := http.NewRequest(GET, fmt.Sprintf("%v/engine/url/%v", apiUrl, id), nil)
		if err != nil {
			// TODO error page
			panic(err.Error())
		}

		request.Header.Set("X-Api-Key", apiKeyA4)

		response, err := client.Do(request)
		if err != nil {
			// TODO error page
			panic(err.Error())
		}
		defer response.Body.Close()

		body := &etools.Body{}
		err = json.NewDecoder(response.Body).Decode(body)
		if err != nil {
			// TODO error page
			panic(err.Error())
		}

		if body.Code == "DBI404" {
			// TODO not found destination
			w.Header().Set("Location", "https://google.com")
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		if body.Code != "S200" {
			// TODO error page
			w.Header().Set("Location", "https://google.com")
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		data := body.Data.(map[string]any)

		destination, ok := data["destination"]
		destinationstring, ok := destination.(string)
		if !ok {
			// TODO error page
			w.Header().Set("Location", "https://google.com")
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		fmt.Printf("%+v\n", body)
		w.Header().Set("Location", destinationstring)
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	http.HandleFunc(buildPath("/url", &POST), func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(etools.NewBody(nil,
				"content type not supported", etools.CODE_BAD_REQUEST),
			)
			return
		}
		type UrlMicro struct {
			Destination string `dynamodbav:"destination" json:"destination"`
		}

		murl := &UrlMicro{}
		err := json.NewDecoder(r.Body).Decode(murl)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(etools.NewBody(nil,
				"invalid json", etools.CODE_BAD_REQUEST),
			)
			return
		}

		now := time.Now().Unix()
		thirtyDays := 30 * (time.Hour.Seconds() * 24)
		url := &types.UrlBase{
			Rash:        "p-sadw",
			Destination: murl.Destination,
			// Public url last 30 days
			Ttl: int(now) + int(thirtyDays),
		}

		var byt bytes.Buffer
		err = json.NewEncoder(&byt).Encode(&url)
		if err != nil {
			panic(err.Error())
		}
		request, err := http.NewRequest(POST, fmt.Sprintf("%v/engine/url", apiUrl), &byt)
		if err != nil {
			// TODO error page
			panic(err.Error())
		}

		request.Header.Set("X-Api-Key", apiKeyA4)
		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		if err != nil {
			// TODO error page
			panic(err.Error())
		}
		defer response.Body.Close()

		body := &etools.Body{}
		err = json.NewDecoder(response.Body).Decode(body)
		if err != nil {
			// TODO error page
			panic(err.Error())
		}

		fmt.Printf("%+v\n", body)
		if body.Code != "S200" {
			// TODO error page
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(etools.NewBody(nil,
				"Internal server errro", etools.CODE_INTERNAL_SERVER_ERROR),
			)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(etools.NewBody(nil, "Ok", etools.CODE_OK))
	})

	httpLambda = httpadapter.New(http.DefaultServeMux)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	httpLambda.StripBasePath("/shortener/v1")
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
