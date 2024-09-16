package tools

import (
	"fmt"
	"os"
)

type Consts map[string]string

const DEFAULT_ERROR_PATH = "oops"

var GET_DEFAULT_WEBURL = func(dev bool) string {
	if dev {
		return "https://s-dev.zenhalab.com"
	}
	return "https://s.zenhalab.com"
}

var ProdConsts = Consts{
	"DEFAULT_ERROR_PAGE":     fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(false), DEFAULT_ERROR_PATH),
	"DEFAULT_NOT_FOUND_PAGE": fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(false), "404"),
	"DEFAULT_WEBURL":         GET_DEFAULT_WEBURL(false),
	"DEFAULT_ERROR_PATH":     "oops",
}

var DevConsts = Consts{
	"DEFAULT_ERROR_PAGE":     fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(true), DEFAULT_ERROR_PATH),
	"DEFAULT_NOT_FOUND_PAGE": fmt.Sprintf("%v/%v", GET_DEFAULT_WEBURL(true), "404"),
	"DEFAULT_WEBURL":         GET_DEFAULT_WEBURL(true),
	"DEFAULT_ERROR_PATH":     "oops",
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
