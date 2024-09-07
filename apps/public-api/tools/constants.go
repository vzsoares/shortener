package tools

import "fmt"

const DEFAULT_ERROR_PATH = "oops"
const DEFAULT_WEBURL = "https://s.zenhalab.com"

var DEFAULT_ERROR_PAGE = fmt.Sprintf("%v/%v", DEFAULT_WEBURL, DEFAULT_ERROR_PATH)
var DEFAULT_NOT_FOUND_PAGE = fmt.Sprintf("%v/%v", DEFAULT_WEBURL, "404")
