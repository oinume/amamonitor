package fetcher

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func Test_Fetcher_GiftissueClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		f, err := os.Open("testdata/giftissue.html")
		if err != nil {
			t.Fatalf("os.Open failed: %v", err)
		}
		defer func() { _ = f.Close() }()
		if _, err := io.Copy(w, f); err != nil {
			t.Fatalf("io.Copy failed: %v", err)
		}
	}))
	defer ts.Close()

	f := NewWithClient(NewGiftissueClient())
	giftItems, err := f.Fetch(context.Background(), &FetchOptions{
		URL: ts.URL,
	})
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	want := []*GiftItem{
		NewGiftItem(GiftissueProvider, "89.8", 20000, 17960),
		NewGiftItem(GiftissueProvider, "89.8", 25000, 22450),
	}
	if len(giftItems) != len(want) {
		t.Fatalf("unexpected Fetch result giftItems length")
	}
	if got := giftItems; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected Fetch result giftItems: got=%+v, want=%+v", got, want)
	}
}
