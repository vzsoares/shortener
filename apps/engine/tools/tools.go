package tools

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type Body struct {
	data    any
	message string
	code    CODES
}

func GatewayResponse(statusCode int, data any, msg string, code CODES) events.APIGatewayV2HTTPResponse {
	var marshalled []byte
	var err error
	var body Body = Body{
		data:    data,
		message: msg,
		code:    code,
	}

	marshalled, err = json.Marshal(body)
	if err != nil {
		return GatewayResponse(http.StatusInternalServerError, nil, err.Error(), CODE_INTERNAL_SERVER_ERROR)
	}

	stringified := string(marshalled)
	if stringified == "null" {
		stringified = ""
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            stringified,
		IsBase64Encoded: false,
	}
}
