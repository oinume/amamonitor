package fetcher

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_AmatenClient_Fetch(t *testing.T) {
	// TODO: httptest.NewServer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		gifts := []AmatenGift{
			{
				ID:        123,
				FaceValue: 10000,
				Price:     8710,
				Rate:      "87.1",
			},
			{
				ID:        456,
				FaceValue: 1000,
				Price:     900,
				Rate:      "90.0",
			},
		}
		response := AmatenGiftResponse{
			Gifts: gifts,
		}

		if err := json.NewEncoder(w).Encode(&response); err != nil {
			t.Fatalf("Encode() failed: %v", err)
		}
	}))
	defer ts.Close()

	c, _ := NewAmatenClient()
	giftItems, err := c.Fetch(context.Background(), &FetchOptions{
		URL: ts.URL,
	})
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	want := []*GiftItem{
		NewGiftItem(AmatenProvider, "87.1", 10000, 8710),
		NewGiftItem(AmatenProvider, "90.0", 1000, 900),
	}
	if len(giftItems) != len(want) {
		t.Fatalf("unexpected giftItems length")
	}
	if got := giftItems; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected Fetch result: got=%+v, want=%+v", got, want)
	}
}

//func Test_amatenClient_FetchHTML(t *testing.T) {
//	c, _ := NewAmatenClient()
//	html, err := c.FetchHTML(context.Background(), fetchURL)
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Printf("--- HTML ---\n%s\n", html)
//}
/*
{
  "user_signed_in": false,
  "gifts": [
    {
      "id": 3840869,
      "revision": 1,
      "f
*/
