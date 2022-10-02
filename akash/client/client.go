package client

import (
	"context"
)

type AkashClient struct {
	ctx             context.Context
	Config          AkashConfiguration
	transactionNote string
}

type AkashConfiguration struct {
	KeyName        string
	KeyringBackend string
	AccountAddress string
	Net            string
	Version        string
	ChainId        string
	Node           string
	Home           string
	Path           string
}

func (ak *AkashClient) GetContext() context.Context {
	return ak.ctx
}

func (ak *AkashClient) GetPath() string {
	return ak.Config.Path
}

func (ak *AkashClient) SetGlobalTransactionNote(note string) {
	ak.transactionNote = note
}

func New(ctx context.Context, configuration AkashConfiguration) *AkashClient {
	return &AkashClient{ctx: ctx, Config: configuration}
}
