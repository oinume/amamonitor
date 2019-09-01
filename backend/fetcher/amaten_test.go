package fetcher

import (
	"context"
	"fmt"
	"testing"
)

func Test_amatenClient_Fetch(t *testing.T) {
	c, _ := NewAmatenClient()
	_, err := c.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func Test_amatenClient_FetchHTML(t *testing.T) {
	c, _ := NewAmatenClient()
	html, err := c.FetchHTML(context.Background(), fetchURL)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("--- HTML ---\n%s\n", html)
}
