package tools

import (
	"fmt"
	"libs/utils"
	"os"
)

const DEFAULT_ERROR_PATH = "oops.html"
const DEFAULT_NOT_FOUND_PATH = "404.html"

var BASE_API_URL_LOCAL string = "http://localhost:4000"

var GET_API_BASE_URL = func() string {
	if DEBUG {
		return BASE_API_URL_LOCAL
	}
	return fmt.Sprintf("%v/%v", os.Getenv("API_BASE_URL"), "/shortener/v1")
}

var ConstsMap = utils.ConstsMap{
	"DEFAULT_ERROR_PAGE":     fmt.Sprintf("%v/%v", os.Getenv("FRONT_BASE_URL"), DEFAULT_ERROR_PATH),
	"DEFAULT_NOT_FOUND_PAGE": fmt.Sprintf("%v/%v", os.Getenv("FRONT_BASE_URL"), DEFAULT_NOT_FOUND_PATH),
	"DEFAULT_WEBURL":         os.Getenv("FRONT_BASE_URL"),
	"DEFAULT_ERROR_PATH":     DEFAULT_ERROR_PATH,
	"API_BASE_URL":           GET_API_BASE_URL(),
}

var Consts = utils.NewConsts(os.Getenv("STAGE"), ConstsMap)
