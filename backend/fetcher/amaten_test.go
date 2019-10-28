package fetcher

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_Fetcher_AmatenClient(t *testing.T) {
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

	fetcher := NewWithClient(NewAmatenClient())
	giftItems, err := fetcher.Fetch(context.Background(), &FetchOptions{
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
		t.Fatalf("unexpected Fetch result giftItems length")
	}
	if got := giftItems; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected Fetch result giftItems: got=%+v, want=%+v", got, want)
	}
}
