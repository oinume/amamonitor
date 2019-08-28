// Command click is a chromedp example demonstrating how to use a selector to
// click on an element.
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/oinume/amamonitor/cli"
)

func main() {
	m := &fetcherMain{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type fetcherMain struct {
	outStream io.Writer
	errStream io.Writer
}

func (m *fetcherMain) run(args []string) error {
	flagSet := flag.NewFlagSet("fetcher", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
		server = flagSet.Bool("server", false, "Run as a HTTP server")
		//concurrency     = flagSet.Int("concurrency", 1, "Concurrency of crawler. (default=1)")
		//continueOnError = flagSet.Bool("continue", true, "Continue to crawl if any error occurred. (default=true)")
		//specifiedIDs    = flagSet.String("ids", "", "Teacher IDs")
		//followedOnly    = flagSet.Bool("followedOnly", false, "Crawl followedOnly teachers")
		//all             = flagSet.Bool("all", false, "Crawl all teachers ordered by evaluation")
		//newOnly         = flagSet.Bool("new", false, "Crawl all teachers ordered by new")
		//interval        = flagSet.Duration("interval", 1*time.Second, "Fetch interval. (default=1s)")
		//logLevel        = flag.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	if *server {
		port := os.Getenv("PORT")
		return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}

	return nil
}

/*
func main() {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	err := chromedp.Run(ctx,
		//chromedp.ExecPath()
		chromedp.Navigate(`https://golang.org/pkg/time/`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// find and click "Expand All" link
		chromedp.Click(`#pkg-examples > div`, chromedp.NodeVisible),
		// retrieve the value of the textarea
		chromedp.Value(`#example_After .play .input textarea`, &example),
	)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	log.Printf("Go's time.After example:\n%s", example)
}
*/
