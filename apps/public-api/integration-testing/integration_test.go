package main

import (
	"apps/public-api/handlers"
	"context"
	"encoding/json"
	"libs/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
)

var apiUrl string
var apiKeyA4 string
var parameterStore *utils.Ssm

func init() {
	apiUrl = "http://localhost:4000"

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

}

func Test_GetHandler_RedirectError(t *testing.T) {
	ctx := context.TODO()
	client := http.Client{}

	handler := handlers.NewHttpHandler(ctx, client, apiUrl, apiKeyA4)

	// Create Request
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	// Exec Request
	handler.GetHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	// Parse res body
	target := &utils.Body{}
	err := json.NewDecoder(res.Body).Decode(target)

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if target.Code != "DBI404" {
		t.Errorf("expected DBI404 got %v", target.Code)
	}
}
