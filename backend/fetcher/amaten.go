package fetcher

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const (
	amatenTargetURL = "https://amaten.com/api/gifts?order=&type=amazon&limit=20&last_id="
)

type amatenGiftResponse struct {
	Gifts []amatenGift `json:"gifts"`
}

type amatenGift struct {
	ID        int    `json:"id"`
	FaceValue uint   `json:"face_value"`
	Price     uint   `json:"price"`
	Rate      string `json:"rate"`
}

func NewAmatenClient() (*amatenClient, error) {
	return &amatenClient{
		httpClient: GetDefaultHTTPClient(),
	}, nil
}

type amatenClient struct {
	httpClient *http.Client
}

func (c *amatenClient) Fetch(ctx context.Context, options *FetchOptions) ([]*GiftItem, error) {
	targetURL := amatenTargetURL
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
	return c.decodeJSON(resp.Body)
}

func (c *amatenClient) setHeaders(req *http.Request) {
	headers := map[string]string{
		"User-Agent":       UserAgent,
		"Sec-Fetch-Mode":   "cors",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://amaten.com/exhibitions/amazon",
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}

func (c *amatenClient) decodeJSON(reader io.Reader) ([]*GiftItem, error) {
	/* response
	   {
	     "id": 3834718,
	     "revision": 0,
	     "face_value": 10000,
	     "price": 8710,
	     "type": "amazon",
	     "rate": "87.1",
	     "is_mine": false,
	     "cnt": 5,
	     "users_total_count": 297282,
	     "users_error_count": 1184
	   }
	*/
	var r amatenGiftResponse
	if err := json.NewDecoder(reader).Decode(&r); err != nil {
		return nil, err
	}
	//fmt.Printf("r.gifts = %+v\n", r.Gifts)

	giftItems := make([]*GiftItem, len(r.Gifts))
	for i, gift := range r.Gifts {
		giftItems[i] = NewGiftItem(gift.Rate, gift.Price)
	}
	return giftItems, nil
}
