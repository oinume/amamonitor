package fetcher

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/oinume/amamonitor/backend/model"
)

type Provider string

const (
	AmatenProvider    Provider = "amaten"
	GiftissueProvider Provider = "giftissue"
	UserAgent                  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
)

func (p Provider) ModelValue() model.Provider {
	switch p {
	case AmatenProvider:
		return model.ProviderAmaten
	case GiftissueProvider:
		return model.ProviderGiftissue
	}
	return model.Provider(0)
}

func NewGiftItem(
	provider Provider,
	discountRate string,
	catalogPrice,
	salesPrice uint,
) *GiftItem {
	return &GiftItem{
		Provider:     provider.ModelValue(),
		DiscountRate: discountRate,
		CatalogPrice: catalogPrice,
		SalesPrice:   salesPrice,
	}
}

type GiftItem struct {
	Provider     model.Provider `json:"provider"`
	DiscountRate string         `json:"discountRate"`
	CatalogPrice uint           `json:"catalogPrice"`
	SalesPrice   uint           `json:"salesPrice"`
}

type FetchOptions struct {
	URL string
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

func GetDefaultHTTPClient() *http.Client {
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

type ProviderClient interface {
	newRequest(options *FetchOptions) (*http.Request, error)
	getHeaders() map[string]string
	parse(body io.Reader) ([]*GiftItem, error)
}

type Fetcher struct {
	client     ProviderClient
	httpClient *http.Client
}

func (f *Fetcher) Fetch(ctx context.Context, options *FetchOptions) ([]*GiftItem, error) {
	req, err := f.client.newRequest(options)
	if err != nil {
		return nil, err
	}
	headers := f.client.getHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return f.client.parse(resp.Body)
}

func NewWithClient(client ProviderClient) *Fetcher {
	return &Fetcher{
		client:     client,
		httpClient: GetDefaultHTTPClient(),
	}
}

func NewFromProvider(p Provider) (*Fetcher, error) {
	var c ProviderClient
	switch p {
	case AmatenProvider:
		c = NewAmatenClient()
	case GiftissueProvider:
		c = NewGiftissueClient()
	default:
		return nil, fmt.Errorf("NewFromProvider failed (unknown provider)")
	}
	return NewWithClient(c), nil
}
