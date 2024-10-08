package tools

import (
	"libs/utils"
	"os"
)

const LOCAL_DB_URL = "http://localhost:8000"
const LOCAL_DB_URL_DOCKER = "http://host.docker.internal:8000"
const REMOTE_DB_URL = "https://dynamodb.us-east-1.amazonaws.com"

var ConstsMap = utils.ConstsMap{
	"URL_TABLE_NAME": os.Getenv("DYNAMO_URL_TABLE_NAME"),
}

var Consts = utils.NewConsts(os.Getenv("STAGE"), ConstsMap)
