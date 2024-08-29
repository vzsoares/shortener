package types

import (
	"context"
)

type UrlStore interface {
	Get(context.Context, string) (*Url, error)
	Put(context.Context, *Url) error
	Delete(context.Context, string) error
}
