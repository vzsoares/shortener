package store

import (
	"apps/engine/types"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableName = "urls"
var BaseEndpoint = "http://localhost:8000/"

type Store interface {
	Get(context.Context, string) (*types.Url, error)
	Put(context.Context, *types.Url) error
	Delete(context.Context, string) error
}

type DynamoStore struct {
	Table  *string
	Client *dynamodb.Client
}

var _ Store = (*DynamoStore)(nil)

func NewDynamoStore(ctx context.Context) *DynamoStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointOptions.DisableHTTPS = true
		o.BaseEndpoint = aws.String(BaseEndpoint)
	})

	return &DynamoStore{
		Table:  &TableName,
		Client: client,
	}
}

func (s *DynamoStore) Get(ctx context.Context, rash string) (*types.Url, error) {
	input := &dynamodb.GetItemInput{
		TableName: &TableName,
		Key: map[string]ddbtypes.AttributeValue{
			"Rash": &ddbtypes.AttributeValueMemberS{Value: rash},
		}}

	res, err := s.Client.GetItem(ctx, input)
	if err != nil {
		panic(err.Error())
	}
	if res.Item == nil {
		panic("404")
	}

	item := types.Url{}
	err = attributevalue.UnmarshalMap(res.Item, &item)
	if err != nil {
		panic(err.Error())
	}

	return &item, nil
}

func (s *DynamoStore) Put(ctx context.Context, url *types.Url) error {
	// Destination
	update := expression.Set(
		expression.Name("Destination"),
		expression.Value(url.Rash),
	)
	// Version + 1
	update.Add(expression.Name("Version"), expression.Value(1))
	// CreatedAt if not set
	update.Set(expression.Name("CreatedAt"),
		expression.IfNotExists(expression.Name("CreatedAt"),
			expression.Value(url.CreatedAt),
		),
	)
	// UpdatedAt
	update.Set(expression.Name("UpdatedAt"), expression.Value(url.UpdatedAt))
	// Ttl
	update.Set(expression.Name("Ttl"), expression.Value(url.Ttl))

	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		panic(err.Error())
	}

	res, err := s.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &TableName,
		Key: map[string]ddbtypes.AttributeValue{
			"Rash": &ddbtypes.AttributeValueMemberS{Value: url.Rash},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              ddbtypes.ReturnValueAllNew,
	})
	if err != nil {
		a := expr.Names()
		b := expr.Values()
		c := expr.Update()

		bs, _ := json.Marshal(a)
		bss, _ := json.Marshal(b)
		bsss, _ := json.Marshal(c)
		fmt.Println(string(bs), string(bss), string(bsss))
		panic(err.Error())
	}
	b := &types.Url{}

	err = attributevalue.UnmarshalMap(res.Attributes, b)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("***** %+v\n", b)
	return err
}

func (s *DynamoStore) Delete(ctx context.Context, rash string) error {
	input := &dynamodb.DeleteItemInput{
		TableName: &TableName,
		Key: map[string]ddbtypes.AttributeValue{
			"Rash": &ddbtypes.AttributeValueMemberS{Value: rash},
		},
	}

	_, err := s.Client.DeleteItem(ctx, input)
	if err != nil {
		panic(err.Error())
	}

	return err
}
