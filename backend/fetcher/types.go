package fetcher

import (
	"context"
	"fmt"
	"net/url"
)

type Type string

const (
	amatenType    Type = "amaten.com"
	giftissueType Type = "giftissue.com"
)

type GiftCard struct {
	discountRate string
	salesPrice   uint
}

func (gc *GiftCard) SalesPrice() uint {
	return gc.salesPrice
}

func (gc *GiftCard) DiscountRate() string {
	return gc.discountRate
}

type Client interface {
	Fetch(ctx context.Context) ([]*GiftCard, error)
	FetchHTML(ctx context.Context, url string) (string, error)
}

func NewClientFromType(t Type) (Client, error) {
	switch t {
	case amatenType:
		return NewAmatenClient()
	case giftissueType:
		return NewGiftissueClient()
	}
	return nil, fmt.Errorf("failed to new client (unknown url)")
}

func NewClientFromURL(urlStr string) (Client, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	return NewClientFromType(Type(u.Host))
}
