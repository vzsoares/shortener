package store

import (
	"apps/engine/tools"
	"apps/engine/types"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var TableName = "shortener-urls"

type DynamoStore struct {
	Table  *string
	Client *dynamodb.Client
}

var _ types.UrlStore = (*DynamoStore)(nil)

func NewDynamoStore(ctx context.Context, endpoint string, https bool) *DynamoStore {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err.Error())
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.EndpointOptions.DisableHTTPS = https
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &DynamoStore{
		Table:  &TableName,
		Client: client,
	}
}

func (s *DynamoStore) Get(ctx context.Context, rash string) (*types.UrlFull, error) {
	input := &dynamodb.GetItemInput{
		TableName: &TableName,
		Key: map[string]ddbtypes.AttributeValue{
			"Rash": &ddbtypes.AttributeValueMemberS{Value: rash},
		}}

	res, err := s.Client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	if res.Item == nil {
		return nil, tools.ItemNotFoundError
	}

	item := types.UrlFull{}
	err = attributevalue.UnmarshalMap(res.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *DynamoStore) Put(ctx context.Context, url *types.UrlFull) error {
	// Destination
	update := expression.Set(
		expression.Name("Destination"),
		expression.Value(url.Destination),
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
		return err
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
		return err
	}

	b := &types.UrlFull{}
	err = attributevalue.UnmarshalMap(res.Attributes, b)
	if err != nil {
		return err
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
		return err
	}

	return err
}
