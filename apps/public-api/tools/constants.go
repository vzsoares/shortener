package tools

import (
	"fmt"
	"libs/utils"
	"os"
)

const DEFAULT_ERROR_PATH = "oops.html"
const DEFAULT_NOT_FOUND_PATH = "404.html"

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

var ProdConsts = utils.ConstsMap{
	"DEFAULT_ERROR_PAGE":     fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(false), DEFAULT_ERROR_PATH),
	"DEFAULT_NOT_FOUND_PAGE": fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(false), DEFAULT_NOT_FOUND_PATH),
	"DEFAULT_WEBURL":         GET_DEFAULT_WEBURL(false),
	"DEFAULT_ERROR_PATH":     DEFAULT_ERROR_PATH,
	"API_BASE_URL":           GET_API_BASE_URL(false),
}

var DevConsts = utils.ConstsMap{
	"DEFAULT_ERROR_PAGE":     fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(true), DEFAULT_ERROR_PATH),
	"DEFAULT_NOT_FOUND_PAGE": fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(true), DEFAULT_NOT_FOUND_PATH),
	"DEFAULT_WEBURL":         GET_DEFAULT_WEBURL(true),
	"DEFAULT_ERROR_PATH":     DEFAULT_ERROR_PATH,
	"API_BASE_URL":           GET_API_BASE_URL(true),
}

var Consts = utils.NewConsts(os.Getenv("STAGE"), ProdConsts, DevConsts)
