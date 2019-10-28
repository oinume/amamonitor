package fetcher

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const (
	amatenTargetURL = "https://amaten.com/api/gifts?order=&type=amazon&limit=20&last_id="
)

type AmatenGiftResponse struct {
	Gifts []AmatenGift `json:"gifts"`
}

type AmatenGift struct {
	ID        int    `json:"id"`
	FaceValue uint   `json:"face_value"`
	Price     uint   `json:"price"`
	Rate      string `json:"rate"`
}

func NewAmatenClient() *amatenClient {
	return &amatenClient{}
}

type amatenClient struct{}

func (c *amatenClient) newRequest(options *FetchOptions) (*http.Request, error) {
	targetURL := amatenTargetURL
	if options != nil && options.URL != "" {
		targetURL = options.URL
	}
	return http.NewRequest("GET", targetURL, nil)
}

func (c *amatenClient) getHeaders() map[string]string {
	return map[string]string{
		"User-Agent":       UserAgent,
		"Sec-Fetch-Mode":   "cors",
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://amaten.com/exhibitions/amazon",
	}
}

func (c *amatenClient) parse(body io.Reader) ([]*GiftItem, error) {
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
	var r AmatenGiftResponse
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return nil, err
	}
	//fmt.Printf("r.gifts = %+v\n", r.Gifts)

	giftItems := make([]*GiftItem, len(r.Gifts))
	for i, gift := range r.Gifts {
		giftItems[i] = NewGiftItem(AmatenProvider, gift.Rate, gift.FaceValue, gift.Price)
	}
	return giftItems, nil

}

func NewFakeAmatenAPIGiftsHandler(t *testing.T, gifts []AmatenGift) func(w http.ResponseWriter, r *http.Request) {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := AmatenGiftResponse{
			Gifts: gifts,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			t.Fatalf("Encode() failed: %v", err)
		}
	}
}

//func FakeAmatenAPIGiftsHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//
//	gifts := []AmatenGift{
//		{
//			ID:        123,
//			FaceValue: 10000,
//			Price:     8710,
//			Rate:      "87.1",
//		},
//		{
//			ID:        456,
//			FaceValue: 1000,
//			Price:     900,
//			Rate:      "90.0",
//		},
//	}
//	response := amatenGiftResponse{
//		Gifts: gifts,
//	}
//
//	if err := json.NewEncoder(w).Encode(&response); err != nil {
//		//t.Fatalf("Encode() failed: %v", err)
//	}
//})
