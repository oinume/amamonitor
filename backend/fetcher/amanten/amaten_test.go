package amanten

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oinume/amamonitor/backend/fetcher"
)

func Test_AmatenClient_Fetch(t *testing.T) {
	// TODO: httptest.NewServer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// TODO: application/json
		// TODO: Use Same response struct in amaten.go
		gifts := []gift{
			gift{},
		}
		response := giftResponse{
			Gifts: gifts,
		}
		if err := json.NewEncoder(w).Encode(&response); err != nil {
			t.Fatalf("encode failed: %v", err)
		}
	}))
	defer ts.Close()

	c, _ := NewAmatenClient()
	giftItems, err := c.Fetch(context.Background(), &fetcher.FetchOptions{
		URL: ts.URL,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(giftItems) == 0 {
		t.Errorf("giftItems length is zero")
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
