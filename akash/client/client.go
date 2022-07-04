package client

import (
	"context"
)

type AkashClient struct {
	ctx context.Context
}

func New(ctx context.Context) *AkashClient {
	return &AkashClient{ctx: ctx}
}
