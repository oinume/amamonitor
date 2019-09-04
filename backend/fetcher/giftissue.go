package fetcher

import "context"

func NewGiftissueClient() (*GiftissueClient, error) {
	return &GiftissueClient{}, nil
}

type GiftissueClient struct{}

func (c *GiftissueClient) Fetch(ctx context.Context, options *FetchOptions) ([]*GiftItem, error) {
	panic("implement me")
}
