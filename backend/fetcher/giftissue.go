package fetcher

import "context"

func NewGiftissueClient() (*giftissueClient, error) {
	return &giftissueClient{}, nil
}

type giftissueClient struct{}

func (c *giftissueClient) Fetch(ctx context.Context) ([]*GiftCard, error) {
	panic("implement me")
}
