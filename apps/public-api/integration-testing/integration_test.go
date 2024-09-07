package main

import (
	"apps/public-api/handlers"
	"apps/public-api/tools"
	"context"
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

func Test_GetHandler_RedirectNotFound(t *testing.T) {
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

	location := res.Header.Get("Location")

	if res.StatusCode != 307 {
		t.Errorf("expected StatusCode to be 307 got %v", res.StatusCode)
	}

	if location != tools.DEFAULT_NOT_FOUND_PAGE {
		t.Errorf("expected to be %v got %v", tools.DEFAULT_NOT_FOUND_PAGE, location)
	}
}

func Test_GetHandler_RedirectError(t *testing.T) {
	ctx := context.TODO()
	client := http.Client{}

	handler := handlers.NewHttpHandler(ctx, client, apiUrl, apiKeyA4)

	// Create Request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Exec Request
	handler.GetHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	location := res.Header.Get("Location")

	if res.StatusCode != 307 {
		t.Errorf("expected StatusCode to be 307 got %v", res.StatusCode)
	}

	if location != tools.DEFAULT_ERROR_PAGE {
		t.Errorf("expected to be %v got %v", tools.DEFAULT_ERROR_PAGE, location)
	}
}

func Test_PostHandler_BadRequestError(t *testing.T) {
	ctx := context.TODO()
	client := http.Client{}

	handler := handlers.NewHttpHandler(ctx, client, apiUrl, apiKeyA4)

	// Create Request
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	// Exec Request
	handler.PostHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != 400 {
		t.Errorf("expected StatusCode to be 400 got %v", res.StatusCode)
	}
}
