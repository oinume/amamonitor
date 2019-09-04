package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/oinume/amamonitor/backend/cli"
	"github.com/oinume/amamonitor/backend/fetcher"
	"github.com/oinume/amamonitor/backend/http_server"
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
		//logLevel        = flag.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	if *server {
		port := os.Getenv("PORT")
		server := http_server.New()
		return http.ListenAndServe(fmt.Sprintf(":%s", port), server.NewRouter())
	}

	client, err := fetcher.NewClientFromURL("https://amaten.com")
	if err != nil {
		return err
	}
	giftCards, err := client.Fetch(context.Background(), &fetcher.FetchOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("discountRate = %v\n", giftCards[0].DiscountRate())
	fmt.Printf("salesPrice = %v\n", giftCards[0].SalesPrice())

	return nil
}
