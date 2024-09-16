package tools

import (
	"fmt"
	"os"
)

type Consts map[string]string

const DEFAULT_ERROR_PATH = "oops"

var BASE_API_URL_REMOTE_DEV string = "https://api-dev.zenhalab.com/shortener/v1"
var BASE_API_URL_REMOTE_PROD string = "https://api.zenhalab.com/shortener/v1"
var BASE_API_URL_LOCAL string = "http://localhost:4000"
var GET_DEFAULT_WEBURL = func(dev bool) string {
	if dev {
		return "https://s-dev.zenhalab.com"
	}
	return "https://s.zenhalab.com"
}
var GET_API_BASE_URL = func(dev bool) string {
	if DEBUG {
		return BASE_API_URL_LOCAL
	}
	if dev {
		return BASE_API_URL_REMOTE_DEV
	}
	return BASE_API_URL_REMOTE_PROD
}

var ProdConsts = Consts{
	"DEFAULT_ERROR_PAGE":     fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(false), DEFAULT_ERROR_PATH),
	"DEFAULT_NOT_FOUND_PAGE": fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(false), "404"),
	"DEFAULT_WEBURL":         GET_DEFAULT_WEBURL(false),
	"DEFAULT_ERROR_PATH":     "oops",
	"API_BASE_URL":           GET_API_BASE_URL(false),
}

var DevConsts = Consts{
	"DEFAULT_ERROR_PAGE":     fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(true), DEFAULT_ERROR_PATH),
	"DEFAULT_NOT_FOUND_PAGE": fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(true), "404"),
	"DEFAULT_WEBURL":         GET_DEFAULT_WEBURL(true),
	"DEFAULT_ERROR_PATH":     "oops",
	"API_BASE_URL":           GET_API_BASE_URL(true),
}

func GetConst(key string) string {
	var stage, ok = os.LookupEnv("STAGE")
	if !ok {
		panic("No STAGE set")
	}

	var v string
	if stage == "dev" {
		v, ok = DevConsts[key]
	} else {
		v, ok = ProdConsts[key]
	}

	if !ok {
		panic(fmt.Sprintf("Variable not set: %v", key))
	}
	return v
}
