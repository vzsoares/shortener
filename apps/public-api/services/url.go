package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	etools "apps/engine/tools"
	"apps/engine/types"
)

func GetUrl(id string, apiUrl string, apiKeyA4 string, client http.Client) (*etools.Body, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%v/engine/url/%v", apiUrl, id), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Api-Key", apiKeyA4)

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body := &etools.Body{}
	err = json.NewDecoder(response.Body).Decode(body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func PutUrl(url *types.UrlBase, apiUrl string, apiKeyA4 string, client http.Client) (*etools.Body, error) {
	var byt bytes.Buffer
	err := json.NewEncoder(&byt).Encode(url)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%v/engine/url", apiUrl), &byt)
	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Api-Key", apiKeyA4)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body := &etools.Body{}
	err = json.NewDecoder(response.Body).Decode(body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
