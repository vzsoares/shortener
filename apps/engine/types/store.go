package types

import (
	"context"
)

type UrlStore interface {
	Get(context.Context, string) (*UrlFull, error)
	Put(context.Context, *UrlFull) error
	Delete(context.Context, string) error
}
