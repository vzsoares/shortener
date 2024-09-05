package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	etools "apps/engine/tools"
	"apps/engine/types"
	"apps/public-api/services"
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
var errorPageUrl = "https://google.com"
var notFoundPageUrl = "https://google.com"

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
			respondRedirect(w, errorPageUrl)
		}
		body, err := services.GetUrl(id, apiUrl, apiKeyA4, client)
		if err != nil {
			respondRedirect(w, errorPageUrl)
		}

		if body.Code == "DBI404" {
			respondRedirect(w, notFoundPageUrl)
			return
		}
		if body.Code != "S200" {
			respondRedirect(w, errorPageUrl)
			return
		}

		data := body.Data.(map[string]any)
		destination, ok := data["destination"]
		destinationstring, ok := destination.(string)
		if !ok {
			respondRedirect(w, errorPageUrl)
			return
		}

		fmt.Printf("%+v\n", body)
		respondRedirect(w, destinationstring)
	})

	http.HandleFunc(buildPath("/url", &POST), func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			respondJson(w, http.StatusBadRequest,
				etools.NewBody(nil,
					"Content type not supported", etools.CODE_BAD_REQUEST),
			)
			return
		}
		type UrlMicro struct {
			Destination string `dynamodbav:"destination" json:"destination"`
		}

		murl := &UrlMicro{}
		err := json.NewDecoder(r.Body).Decode(murl)

		if err != nil {
			respondJson(w, http.StatusBadRequest,
				etools.NewBody(nil,
					"Invalid json", etools.CODE_BAD_REQUEST),
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
			respondJson(w, http.StatusInternalServerError,
				etools.NewBody(nil,
					"Internal server error", etools.CODE_INTERNAL_SERVER_ERROR),
			)
		}
		request, err := http.NewRequest(POST, fmt.Sprintf("%v/engine/url", apiUrl), &byt)
		if err != nil {
			respondJson(w, http.StatusInternalServerError,
				etools.NewBody(nil,
					"Internal server error", etools.CODE_INTERNAL_SERVER_ERROR),
			)
		}

		request.Header.Set("X-Api-Key", apiKeyA4)
		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		if err != nil {
			respondJson(w, http.StatusInternalServerError,
				etools.NewBody(nil,
					"Internal server error", etools.CODE_INTERNAL_SERVER_ERROR),
			)
		}
		defer response.Body.Close()

		body := &etools.Body{}
		err = json.NewDecoder(response.Body).Decode(body)
		if err != nil {
			respondJson(w, http.StatusInternalServerError,
				etools.NewBody(nil,
					"Internal server error", etools.CODE_INTERNAL_SERVER_ERROR),
			)
		}

		fmt.Printf("%+v\n", body)
		if body.Code != "S200" {
			respondJson(w, http.StatusInternalServerError,
				etools.NewBody(nil,
					"Internal server error", etools.CODE_INTERNAL_SERVER_ERROR),
			)
			return
		}
		var frontBaseUrl = "https://s.zenhalab.com"
		type data struct {
			Url string `json:"url"`
		}
		resData := &data{Url: fmt.Sprintf("%v/%v", frontBaseUrl, url.Rash)}
		resBody := etools.NewBody(resData, "Ok", etools.CODE_OK)

		respondJson(w, http.StatusOK, resBody)
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

var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~()'!*:@,;"

func genRash(size int) string {
	l := len(chars)
	r := ""
	for range size {
		rd := rand.Intn(l)
		c := chars[rd]
		r += string(c)
	}
	return r
}

func respondJson(w http.ResponseWriter, s int, j any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)
	json.NewEncoder(w).Encode(j)
}
func respondRedirect(w http.ResponseWriter, l string) {
	w.Header().Set("Location", l)
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
