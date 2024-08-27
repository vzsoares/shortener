package main

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"time"

	store "apps/engine/store/url"
	"apps/engine/types"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
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

func TestFlow(t *testing.T) {
	fk := getRandomProduct()
	ctx := context.TODO()
	store := store.NewDynamoStore(ctx)

	listTablesRes, err := store.Client.ListTables(ctx, &dynamodb.ListTablesInput{})
	if err != nil {
		panic(err.Error())
	}
	hasTable := slices.Contains(listTablesRes.TableNames, "urls")

	if !hasTable {
		created, err := store.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
			AttributeDefinitions: []ddbtypes.AttributeDefinition{
				// {
				// 	AttributeName: aws.String("Destination"),
				// 	AttributeType: ddbtypes.ScalarAttributeTypeS,
				// },
				// {
				// 	AttributeName: aws.String("CreatedAt"),
				// 	AttributeType: ddbtypes.ScalarAttributeTypeN,
				// },
				// {
				// 	AttributeName: aws.String("UpdatedAt"),
				// 	AttributeType: ddbtypes.ScalarAttributeTypeN,
				// },
				// {
				// 	AttributeName: aws.String("Ttl"),
				// 	AttributeType: ddbtypes.ScalarAttributeTypeN,
				// },
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
	r := store.Put(ctx, fk)
	println(r)
	// TODO create table
	// TODO create item
	// TODO get item
	// TODO update item
	// TODO get item
	// TODO delete item
	// TODO get item
}
