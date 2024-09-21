package tools

import (
	"libs/utils"
	"os"
)

const LOCAL_DB_URL = "http://localhost:8000"
const LOCAL_DB_URL_DOCKER = "http://host.docker.internal:8000"
const REMOTE_DB_URL = "https://dynamodb.us-east-1.amazonaws.com"

var ProdConsts = utils.ConstsMap{
	"URL_TABLE_NAME": "shortener-urls-prod",
}

var DevConsts = utils.ConstsMap{
	"URL_TABLE_NAME": "shortener-urls-dev",
}

var Consts = utils.NewConsts(os.Getenv("STAGE"), ProdConsts, DevConsts)
