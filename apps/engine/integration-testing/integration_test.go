package main

import (
	b64 "encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"apps/engine/types"
)

var apiUrl string

func init() {
	// _apiUrl, ok := os.LookupEnv("API_URL")
	// if !ok {
	// 	panic("Can't find API_URL environment variable")
	// }

	// apiUrl = _apiUrl
	apiUrl = "http://localhost:8000"
}

func randomString(length int) string {
	u := time.Now().UnixNano()
	s := fmt.Sprint(u)
	sEnc := b64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("%x", sEnc)[:length]
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
	s := randomString(10)
	println(s)
	// TODO create table
	// TODO create item
	// TODO get item
	// TODO update item
	// TODO get item
	// TODO delete item
	// TODO get item
}
