package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	etools "apps/engine/tools"
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
