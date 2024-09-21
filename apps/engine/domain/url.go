package domain

import (
	"context"
	"errors"
	"libs/utils"
	"time"

	"apps/engine/types"
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

func (d *UrlDomain) GetUrl(ctx context.Context, id string) (*types.UrlFull, error) {
	url, err := d.store.Get(ctx, id)

	if url != nil && url.Rash != id {
		return nil, utils.IdMismatchError
	}

	return url, err
}

func (d *UrlDomain) PutUrl(ctx context.Context, url *types.UrlBase) error {
	now := time.Now().Unix()

	if url.Ttl > 0 && url.Ttl < int(now) {
		return utils.InputValidationError
	}
	if len(url.Rash) < 4 {
		return utils.InputValidationError
	}
	if len(url.Destination) < 4 {
		return utils.InputValidationError
	}

	full := &types.UrlFull{
		UrlBase: url,
		// dynamo will not overwrite this
		CreatedAt: int(now),
		UpdatedAt: int(now),
		// always 0; dynamo handles increase
		Version: 0,
	}

	err := d.store.Put(ctx, full)

	return err
}

func (d *UrlDomain) DeleteUrl(ctx context.Context, rash string) error {
	if rash == "" {
		return errors.New("Id is empty")
	}

	err := d.store.Delete(ctx, rash)

	return err
}
