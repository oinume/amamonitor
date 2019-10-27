package fetcher

import (
	"context"
	"io"
	"net/http"

	"gopkg.in/xmlpath.v2"
)

const (
	giftissueTargetURL = "https://giftissue.com/ja/category/amazonjp/"
)

func NewGiftissueClient() (*giftissueClient, error) {
	return &giftissueClient{
		httpClient: GetDefaultHTTPClient(),
	}, nil
}

type giftissueClient struct {
	httpClient *http.Client
}

func (c *giftissueClient) Fetch(ctx context.Context, options *FetchOptions) ([]*GiftItem, error) {
	targetURL := giftissueTargetURL
	if options.URL != "" {
		targetURL = options.URL
	}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return c.decodeHTML(resp.Body)
	//return c.decodeJSON(resp.Body)
}

func (c *giftissueClient) setHeaders(req *http.Request) {
	headers := map[string]string{
		"User-Agent": UserAgent,
		//"X-Requested-With": "XMLHttpRequest",
		//"Referer":          "https://amaten.com/exhibitions/amazon",
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (c *giftissueClient) decodeHTML(reader io.Reader) ([]*GiftItem, error) {
	_, err := xmlpath.ParseHTML(reader)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
