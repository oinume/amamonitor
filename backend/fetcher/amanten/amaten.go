package amanten

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/oinume/amamonitor/backend/fetcher"
)

const (
	defaultTargetURL = "https://amaten.com/api/gifts?order=&type=amazon&limit=20&last_id="
)

type giftResponse struct {
	Gifts []gift `json:"gifts"`
}

type gift struct {
	ID        int    `json:"id"`
	FaceValue uint   `json:"face_value"`
	Price     uint   `json:"price"`
	Rate      string `json:"rate"`
}

func NewAmatenClient() (*AmatenClient, error) {
	return &AmatenClient{
		httpClient: fetcher.GetDefaultHTTPClient(),
	}, nil
}

type AmatenClient struct {
	httpClient *http.Client
}

func (c *AmatenClient) Fetch(ctx context.Context, options *fetcher.FetchOptions) ([]*fetcher.GiftItem, error) {
	targetURL := defaultTargetURL
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

func (c *AmatenClient) setHeaders(req *http.Request) {
	headers := map[string]string{
		"User-Agent":       fetcher.UserAgent,
		"Sec-Fetch-Mode":   "cors",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://amaten.com/exhibitions/amazon",
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	//curl 'https://amaten.com/api/gifts?order=&type=amazon&limit=20&last_id=' -H 'Pragma: no-cache' -H 'Sec-Fetch-Site: same-origin' -H 'Accept-Encoding: gzip, deflate, br' -H 'X-CSRF-Token: tYt7Vm86sNzY1TU65VwLcWTibMoaK4nmlhMGPhyun/0=' -H 'Accept-Language: en-US,en;q=0.9,ja;q=0.8' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36' -H 'Sec-Fetch-Mode: cors' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Cache-Control: no-cache' -H 'X-Requested-With: XMLHttpRequest' -H 'Cookie: uid=d94f41c0-e002-4282-8e13-14abac774569; _amaten_session=16c3a8ac5de9eda1a1f60de885c5dcd6; _ga=GA1.2.950863381.1567320935; _gid=GA1.2.347252349.1567320935; _gat=1; _fbp=fb.1.1567320938613.1337910132' -H 'Connection: keep-alive' -H 'Referer: https://amaten.com/exhibitions/amazon' --compressed
}

func (c *AmatenClient) decodeJSON(reader io.Reader) ([]*fetcher.GiftItem, error) {
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
	var r giftResponse
	if err := json.NewDecoder(reader).Decode(&r); err != nil {
		return nil, err
	}
	//fmt.Printf("r.gifts = %+v\n", r.Gifts)

	giftItems := make([]*fetcher.GiftItem, len(r.Gifts))
	for i, gift := range r.Gifts {
		giftItems[i] = fetcher.NewGiftItem(gift.Rate, gift.Price)
	}
	return giftItems, nil
}
