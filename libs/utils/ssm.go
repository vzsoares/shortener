package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Ssm struct {
	client     *ssm.Client
	parameters map[string]string
	context    context.Context
}

func NewSmmStore(cfg aws.Config, ctx context.Context) *Ssm {
	ssmClient := ssm.NewFromConfig(cfg)
	cache := make(map[string]string)

	return &Ssm{
		client:     ssmClient,
		parameters: cache,
		context:    ctx,
	}
}

func (s *Ssm) Get(k string) string {
	v, ok := s.parameters[k]
	if !ok {
		stage, ok := os.LookupEnv("STAGE")
		if !ok {
			panic("No STAGE set")
		}
		res, err := s.client.GetParameter(s.context, &ssm.GetParameterInput{
			Name: aws.String(fmt.Sprintf("/%v/%v", stage, k)),
		})
		if err != nil {
			panic(err.Error())
		}
		vl := *res.Parameter.Value
		v = vl
		s.parameters[k] = vl
	}
	return v
}
