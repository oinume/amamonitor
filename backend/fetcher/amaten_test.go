package fetcher

import (
	"context"
	"testing"
)

func Test_amatenClient_Fetch(t *testing.T) {
	c, _ := NewAmatenClient()
	_, err := c.Fetch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
