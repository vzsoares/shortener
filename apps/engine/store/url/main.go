package main

import (
	"apps/engine/types"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableName = "sample"
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
	// TODO conditional put, if exists then update Version,UpdatedAt, ...rest
	item, err := attributevalue.MarshalMap(url)
	if err != nil {
		panic(err.Error())
	}

	_, err = s.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &TableName,
		Item:      item,
	})
	if err != nil {
		panic(err.Error())
	}

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
