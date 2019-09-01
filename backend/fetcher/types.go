package fetcher

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Type string

const (
	amatenType    Type = "amaten.com"
	giftissueType Type = "giftissue.com"
	userAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
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

type FetchOptions struct {
	URL string
}

type Client interface {
	Fetch(ctx context.Context, options *FetchOptions) ([]*GiftCard, error)
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

var redirectErrorFunc = func(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

var defaultHTTPClient = &http.Client{
	Timeout:       5 * time.Second,
	CheckRedirect: redirectErrorFunc,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		Proxy:               http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 1200 * time.Second,
		}).DialContext,
		IdleConnTimeout:     1200 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			ClientSessionCache: tls.NewLRUClientSessionCache(100),
		},
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func getDefaultHTTPClient() *http.Client {
	//if !config.DefaultVars.EnableFetcherHTTP2 {
	//	return defaultHTTPClient
	//}
	//defaultHTTPClient.Transport = &http2.Transport{
	//	TLSClientConfig: &tls.Config{
	//		ClientSessionCache: tls.NewLRUClientSessionCache(100),
	//	},
	//	StrictMaxConcurrentStreams: true,
	//}
	//return defaultHTTPClient
	return defaultHTTPClient
}
