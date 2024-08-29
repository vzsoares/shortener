package main

import (
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
	"apps/engine/handler"
	store "apps/engine/store"
	"apps/engine/tools"
	"apps/engine/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
)

func urlsAreEqual(a *types.Url, b *types.Url) bool {
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

func getRandomProduct() *types.Url {
	return &types.Url{
		Rash:        randomString(10),
		Destination: randomString(10),
		Ttl:         rand.Int(),
		UpdatedAt:   rand.Int(),
		CreatedAt:   rand.Int(),
		Version:     1,
	}
}

var apiUrl string

func init() {
	// _apiUrl, ok := os.LookupEnv("API_URL")
	// if !ok {
	// 	panic("Can't find API_URL environment variable")
	// }

	// apiUrl = _apiUrl
	apiUrl = "http://localhost:8000"
}

func Test_CompleteFlow_CreateGetAlterGetDeleteGet(t *testing.T) {
	fk := getRandomProduct()
	ctx := context.TODO()
	store := store.NewDynamoStore(ctx)

	listTablesRes, err := store.Client.ListTables(ctx,
		&dynamodb.ListTablesInput{},
	)
	if err != nil {
		t.Error(err.Error())
	}
	hasTable := slices.Contains(listTablesRes.TableNames, "urls")

	//create table
	if !hasTable {
		created, err := store.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
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
			TableName:   aws.String(*store.Table),
			BillingMode: "PAY PER REQUEST",
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
	// TODO create item
	fmt.Printf("*****og %+v\n", fk)
	err = store.Put(ctx, fk)
	if err != nil {
		panic(err.Error())
	}
	// TODO get item
	r, err := store.Get(ctx, fk.Rash)
	if err != nil {
		panic(err.Error())
	}
	ok := urlsAreEqual(r, fk)
	if !ok {
		panic("urls not equal")
	}
	// TODO update item
	tm := &types.Url{
		Rash:      fk.Rash,
		CreatedAt: fk.CreatedAt,
		Version:   2,
	}
	fk.Destination = randomString(10)
	fk.UpdatedAt = rand.Int()
	fk.CreatedAt = rand.Int()
	fk.Version = rand.Int()
	fk.Ttl = rand.Int()

	tm.Destination = fk.Destination
	tm.UpdatedAt = fk.UpdatedAt
	tm.Ttl = fk.Ttl

	err = store.Put(ctx, fk)
	if err != nil {
		panic(err.Error())
	}
	// TODO get item
	r, err = store.Get(ctx, fk.Rash)
	if err != nil {
		panic(err.Error())
	}
	ok = urlsAreEqual(r, tm)
	if !ok {
		panic("urls not equal")
	}
	// TODO delete item
	err = store.Delete(ctx, fk.Rash)
	if err != nil {
		t.Error(err.Error())
	}
	// TODO get item
	r, err = store.Get(ctx, fk.Rash)
	if err == nil {
		t.Error("Must error")
	}
	if !errors.Is(err, tools.ItemNotFoundError) {
		t.Error(err.Error())
	}

}

func Test_GetHandler_NonExistentItem_NotFound(t *testing.T) {
	ctx := context.TODO()
	store := store.NewDynamoStore(ctx)
	domain := domain.NewUrlDomain(ctx, store)
	handler := handler.NewHttpHandler(ctx, domain)

	req := httptest.NewRequest(http.MethodGet, "/upper/123", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.GetHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	target := &tools.Body{}
	err := json.NewDecoder(res.Body).Decode(target)

	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if target.Code != "DBI404" {
		t.Errorf("expected ABC got %v", "123id")
	}
}
