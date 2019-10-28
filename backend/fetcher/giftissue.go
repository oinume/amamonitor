package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/xmlpath.v2"
)

const (
	giftissueTargetURL = "https://giftissue.com/ja/category/amazonjp/"
)

var (
	giftListItemXPath = xmlpath.MustCompile(`//table/tbody/tr[@class='giftList_item']`)
	faceValueXPath    = xmlpath.MustCompile(`td[contains(@class, 'giftList_cell-facevalue')]/span`)
	priceXPath        = xmlpath.MustCompile(`td[contains(@class, 'giftList_cell-price')]/span`)
	rateXPath         = xmlpath.MustCompile(`td[contains(@class, 'giftList_rate')]/span`)
)

func NewGiftissueClient() *giftissueClient {
	return &giftissueClient{}
}

type giftissueClient struct{}

func (c *giftissueClient) newRequest(options *FetchOptions) (*http.Request, error) {
	targetURL := giftissueTargetURL
	if options != nil && options.URL != "" {
		targetURL = options.URL
	}
	return http.NewRequest("GET", targetURL, nil)
}

func (c *giftissueClient) getHeaders() map[string]string {
	return map[string]string{
		"User-Agent": UserAgent,
		//"X-Requested-With": "XMLHttpRequest",
		//"Referer":          "https://amaten.com/exhibitions/amazon",
	}
}

func (c *giftissueClient) parse(body io.Reader) ([]*GiftItem, error) {
	root, err := xmlpath.ParseHTML(body)
	if err != nil {
		return nil, err
	}

	giftItems := make([]*GiftItem, 0, 100)
	for iter := giftListItemXPath.Iter(root); iter.Next(); {
		node := iter.Node()
		faceValueString, ok := faceValueXPath.String(node)
		if !ok {
			return nil, fmt.Errorf("failed to parse faceValue with XPath")
		}
		faceValue, err := c.normalizeAmount(faceValueString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse faceValue: %v", err)
		}

		priceString, ok := priceXPath.String(node)
		if !ok {
			return nil, fmt.Errorf("failed to parse price with XPath")
		}
		price, err := c.normalizeAmount(priceString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse price: %v", err)
		}

		rateString, ok := rateXPath.String(node)
		if !ok {
			return nil, fmt.Errorf("failed to parse rate with XPath")
		}
		rate, err := c.normalizeRate(rateString)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rate: %v", err)
		}

		giftItem := NewGiftItem(GiftissueProvider, rate, faceValue, price)
		giftItems = append(giftItems, giftItem)
	}

	return giftItems, nil
}

func (c *giftissueClient) normalizeAmount(s string) (uint, error) {
	s = strings.ReplaceAll(s, "¥", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (c *giftissueClient) normalizeRate(s string) (string, error) {
	v := strings.ReplaceAll(s, "%", "")
	v = strings.TrimSpace(v)
	// TODO: validate
	return v, nil
}
