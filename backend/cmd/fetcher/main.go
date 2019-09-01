package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/oinume/amamonitor/backend/cli"
	"github.com/oinume/amamonitor/backend/fetcher"
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
		return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	}

	client, err := fetcher.NewClientFromURL("https://amaten.com")
	if err != nil {
		return err
	}
	giftCards, err := client.Fetch("https://amaten.com")
	if err != nil {
		return err
	}

	fmt.Printf("discountRate = %v\n", giftCards[0].DiscountRate())
	fmt.Printf("salesPrice = %v\n", giftCards[0].SalesPrice())

	return nil
}
