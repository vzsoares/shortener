package main

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"apps/engine/domain"
	"apps/engine/handlers"
	"apps/engine/store"
	"apps/engine/tools"
	"apps/engine/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func urlsAreEqual(a *types.UrlFull, b *types.UrlFull) bool {
	if a.Rash != b.Rash {
		println("rash")
		return false
	}
	if a.Destination != b.Destination {
		return false
	}
	if a.Ttl != b.Ttl {
		return false
	}
	if a.CreatedAt != b.CreatedAt {
		return false
	}
	if a.UpdatedAt != b.UpdatedAt {
		return false
	}
	if a.Version != b.Version {
		return false
	}
	return true
}

func randomString(length int) string {
	u := time.Now().UnixNano()
	s := fmt.Sprint(u)
	rs := ""
	for i := len(s) - 1; i >= 0; i-- {
		rs = rs + string(s[i])
	}
	sEnc := b64.StdEncoding.EncodeToString([]byte(rs))
	return fmt.Sprintf("%x", sEnc)[:length]
}

func getRandomProduct() *types.UrlFull {
	return &types.UrlFull{
		UrlBase: &types.UrlBase{
			Rash:        randomString(10),
			Destination: randomString(10),
			Ttl:         rand.Int(),
		},
		UpdatedAt: rand.Int(),
		CreatedAt: rand.Int(),
		Version:   1,
	}
}

var apiUrl string
var dstore *store.DynamoStore

func init() {
	// _apiUrl, ok := os.LookupEnv("API_URL")
	// if !ok {
	// 	panic("Can't find API_URL environment variable")
	// }

	// apiUrl = _apiUrl
	apiUrl = tools.LOCAL_DB_URL

	ctx := context.TODO()
	dstore = store.NewDynamoStore(ctx, apiUrl, true)

	listTablesRes, err := dstore.Client.ListTables(ctx,
		&dynamodb.ListTablesInput{},
	)
	if err != nil {
		panic(err.Error())
	}
	hasTable := slices.Contains(listTablesRes.TableNames, *dstore.Table)

	//create table
	if !hasTable {
		created, err := dstore.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
			AttributeDefinitions: []ddbtypes.AttributeDefinition{
				{
					AttributeName: aws.String("Rash"),
					AttributeType: ddbtypes.ScalarAttributeTypeS,
				},
			},
			KeySchema: []ddbtypes.KeySchemaElement{{
				AttributeName: aws.String("Rash"),
				KeyType:       ddbtypes.KeyTypeHash,
			}},
			TableName: dstore.Table,
			ProvisionedThroughput: &ddbtypes.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
		})
		if err != nil {
			panic(err.Error())
		}
		println("created: ", *created.TableDescription.TableName)
	}
}

func Test_StoreCompleteFlow_CreateGetAlterGetDeleteGet(t *testing.T) {
	// Create
	ctx := context.TODO()
	fk := getRandomProduct()
	err := dstore.Put(ctx, fk)
	if err != nil {
		panic(err.Error())
	}

	// Get created
	r, err := dstore.Get(ctx, fk.Rash)
	if err != nil {
		panic(err.Error())
	}
	ok := urlsAreEqual(r, fk)
	if !ok {
		panic("urls not equal")
	}

	// Update created
	tm := &types.UrlFull{
		// Rash, CreatedAt must not change
		UrlBase: &types.UrlBase{Rash: fk.Rash},

		CreatedAt: fk.CreatedAt,
		// Version goes up 1 by 1 automatically
		Version: 2,
	}
	fk.Destination = randomString(10)
	fk.UpdatedAt = rand.Int()
	fk.CreatedAt = rand.Int()
	fk.Version = rand.Int()
	fk.Ttl = rand.Int()

	tm.Destination = fk.Destination
	tm.UpdatedAt = fk.UpdatedAt
	tm.Ttl = fk.Ttl

	err = dstore.Put(ctx, fk)
	if err != nil {
		panic(err.Error())
	}

	// Get updated
	r, err = dstore.Get(ctx, fk.Rash)
	if err != nil {
		panic(err.Error())
	}
	ok = urlsAreEqual(r, tm)
	if !ok {
		panic("urls not equal")
	}

	// Delete item
	err = dstore.Delete(ctx, fk.Rash)
	if err != nil {
		t.Error(err.Error())
	}

	// Get deleted item
	r, err = dstore.Get(ctx, fk.Rash)
	if err == nil {
		t.Error("Must error")
	}
	if !errors.Is(err, tools.ItemNotFoundError) {
		t.Error(err.Error())
	}

}

func Test_GetHandler_NonExistentItem_NotFound(t *testing.T) {
	// Setup
	ctx := context.TODO()
	domain := domain.NewUrlDomain(ctx, dstore)
	handler := handlers.NewHttpHandler(ctx, domain)

	// Create Request
	req := httptest.NewRequest(http.MethodGet, "/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	// Exec Request
	handler.GetHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	// Parse res body
	target := &tools.Body{}
	err := json.NewDecoder(res.Body).Decode(target)

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if target.Code != "DBI404" {
		t.Errorf("expected DBI404 got %v", target.Code)
	}
}

func Test_PostHandler_JustCreate(t *testing.T) {
	// Setup
	ctx := context.TODO()
	domain := domain.NewUrlDomain(ctx, dstore)
	handler := handlers.NewHttpHandler(ctx, domain)

	url := types.UrlBase{
		Rash:        randomString(10),
		Destination: randomString(10),
		Ttl:         0,
	}
	var byt bytes.Buffer
	err := json.NewEncoder(&byt).Encode(&url)
	if err != nil {
		t.Error("failed to parse json")
	}

	// Create Request
	req := httptest.NewRequest(http.MethodPost, "/", &byt)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Exec Request
	handler.PostHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	// Parse res body
	target := &tools.Body{}
	err = json.NewDecoder(res.Body).Decode(target)

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
}
