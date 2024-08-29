package domain

import (
	"apps/engine/types"
	"context"
)

type UrlDomain struct {
	store types.UrlStore
	ctx   context.Context
}

func NewUrlDomain(ctx context.Context, s types.UrlStore) *UrlDomain {
	return &UrlDomain{
		store: s,
		ctx:   ctx,
	}
}

func (d *UrlDomain) GetUrl(ctx context.Context, id string) (*types.Url, error) {

	return d.store.Get(ctx, id)
}
