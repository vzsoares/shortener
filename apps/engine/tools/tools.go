package tools

import (
	"encoding/json"
	"libs/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func GatewayResponse(statusCode int, data any, msg string, code utils.CODES) events.APIGatewayV2HTTPResponse {
	var marshalled []byte
	var err error
	var body utils.Body = utils.Body{
		Data:    data,
		Message: msg,
		Code:    code,
	}

	marshalled, err = json.Marshal(body)
	if err != nil {
		return GatewayResponse(http.StatusInternalServerError, nil, err.Error(), utils.CODE_INTERNAL_SERVER_ERROR)
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
