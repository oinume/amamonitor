package fetcher

import (
	"context"
	"testing"
)

func Test_AmatenClient_Fetch(t *testing.T) {
	// TODO: httptest.NewServer
	c, _ := NewAmatenClient()
	giftCards, err := c.Fetch(context.Background(), &FetchOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(giftCards) == 0 {
		t.Errorf("giftCards length is zero")
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
