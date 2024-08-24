package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"apps/engine/types"
)

var apiUrl string

func init() {
	_apiUrl, ok := os.LookupEnv("API_URL")
	if !ok {
		panic("Can't find API_URL environment variable")
	}

	apiUrl = _apiUrl
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func getRandomProduct() types.Url {
	return types.Url{
		Rash:        randomString(10),
		Destination: randomString(10),
		Ttl:         rand.Int(),
		UpdatedAt:   rand.Int(),
		CreatedAt:   rand.Int(),
	}
}

func TestFlow(t *testing.T) {
	// TODO create table
	// TODO create item
	// TODO get item
	// TODO update item
	// TODO get item
	// TODO delete item
	// TODO get item
}
