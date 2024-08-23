package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHello(t *testing.T) {
	result, err := Handler(context.TODO(), events.APIGatewayV2HTTPRequest{
		Version: "123",
		Body:    "",
	})
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("t: %v\n", result)
}
