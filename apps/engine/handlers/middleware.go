package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var A4Key string

func GetApiKey(key string) string {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}
	ssmClient := ssm.NewFromConfig(cfg)
	res, err := ssmClient.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(key),
	})
	if err != nil {
		panic(err.Error())
	}

	return *res.Parameter.Value
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		incomingKey := r.Header.Get("X-Api-Key")
		if incomingKey == "" {
			http.Error(w, "\"Missing X-Api-Key\"", http.StatusForbidden)
			return
		}

		if A4Key == "" {
			A4Key = GetApiKey("/dev/API_KEY_A4")
		}

		if incomingKey != A4Key {
			http.Error(w, "\"Unauthorized\"", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
