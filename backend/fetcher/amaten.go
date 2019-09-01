package fetcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

const (
	fetchURL = "https://amaten.com/exhibitions/amazon"
)

func NewAmatenClient() (*amatenClient, error) {
	return &amatenClient{}, nil
}

type amatenClient struct {
	chromeDpContext context.Context
}

func (c *amatenClient) Fetch(ctx context.Context) ([]*GiftCard, error) {
	// create chrome instance
	chromeDpCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	c.chromeDpContext = chromeDpCtx
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var title string
	err := chromedp.Run(
		c.chromeDpContext,
		chromedp.Navigate(fetchURL),
		// wait for footer element is visible (ie, page is loaded)
		//chromedp.WaitVisible(`body > footer`),
		//// find and click "Expand All" link
		//chromedp.Click(`#pkg-examples > div`, chromedp.NodeVisible),
		chromedp.Text(`//title`, &title),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	// TODO: Get HTML: https://github.com/chromedp/chromedp/issues/128#issuecomment-497974854
	fmt.Printf("title = %v\n", title)
	return nil, nil
}

func (c *amatenClient) FetchHTML(ctx context.Context, url string) (string, error) {
	// create chrome instance
	chromeDpCtx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	c.chromeDpContext = chromeDpCtx
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(chromeDpCtx, 10*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var html string
	err := chromedp.Run(
		c.chromeDpContext,
		chromedp.Navigate(fetchURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}))
	if err != nil {
		return "", err
	}

	return html, nil
}
