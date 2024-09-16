package handlers

import (
	"apps/engine/types"
	"apps/public-api/services"
	"apps/public-api/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"libs/utils"
	"net/http"
	"time"
)

type UrlHttpHandler struct {
	ctx      context.Context
	client   http.Client
	apiUrl   string
	apiKeyA4 string
}

func NewHttpHandler(ctx context.Context, client http.Client, apiUrl string, apiKeyA4 string) *UrlHttpHandler {
	return &UrlHttpHandler{
		ctx:      ctx,
		client:   client,
		apiUrl:   apiUrl,
		apiKeyA4: apiKeyA4,
	}
}

var errorPageUrl = tools.GetConst("DEFAULT_ERROR_PAGE")
var notFoundPageUrl = tools.GetConst("DEFAULT_NOT_FOUND_PAGE")

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

func getValidRash(client http.Client, size int, apiUrl string, apiKeyA4 string) (string, error) {
	rash := utils.GenUriSafeRash(size, "p-")
	res, err := services.GetUrl(rash, apiUrl, apiKeyA4, client)
	if err != nil {
		return "", err
	}

	if res.Code != "DBI404" {
		return "", errors.New("Exists")
	}
	return rash, nil
}

func (h *UrlHttpHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		respondRedirect(w, errorPageUrl)
	}
	body, err := services.GetUrl(id, h.apiUrl, h.apiKeyA4, h.client)
	if err != nil {
		respondRedirect(w, errorPageUrl)
		return
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
}

func (h *UrlHttpHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		respondJson(w, http.StatusBadRequest,
			utils.NewBody(nil,
				"Content type not supported", utils.CODE_BAD_REQUEST),
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
			utils.NewBody(nil,
				"Invalid json", utils.CODE_BAD_REQUEST),
		)
		return
	}

	now := time.Now().Unix()
	thirtyDays := 30 * (time.Hour.Seconds() * 24)
	fn := func(i int) (string, error) {
		rash, err := getValidRash(h.client, i+5, h.apiUrl, h.apiKeyA4)
		return rash, err
	}
	rash, err := utils.RetryN(fn, 3)

	if err != nil {
		respondJson(w, http.StatusInternalServerError,
			utils.NewBody(nil,
				"Internal server error", utils.CODE_INTERNAL_SERVER_ERROR),
		)
		return
	}
	url := &types.UrlBase{
		Rash:        rash,
		Destination: murl.Destination,
		// Public url last 30 days
		Ttl: int(now) + int(thirtyDays),
	}

	body, err := services.PutUrl(url, h.apiUrl, h.apiKeyA4, h.client)
	if err != nil {
		respondJson(w, http.StatusInternalServerError,
			utils.NewBody(nil,
				"Internal server error", utils.CODE_INTERNAL_SERVER_ERROR),
		)
		return
	}

	if body.Code != "S200" {
		respondJson(w, http.StatusInternalServerError,
			utils.NewBody(nil,
				"Internal server error", utils.CODE_INTERNAL_SERVER_ERROR),
		)
		return
	}
	var frontBaseUrl = "https://s.zenhalab.com"
	type data struct {
		Url string `json:"url"`
	}
	resData := &data{Url: fmt.Sprintf("%v/%v", frontBaseUrl, url.Rash)}
	resBody := utils.NewBody(resData, "Ok", utils.CODE_OK)

	respondJson(w, http.StatusOK, resBody)
}
