package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oinume/amamonitor/backend/cli"
	"github.com/oinume/amamonitor/backend/config"
	"github.com/oinume/amamonitor/backend/fetcher"
	"github.com/oinume/amamonitor/backend/http_server"
	"github.com/oinume/amamonitor/backend/service"
	"github.com/xo/dburl"
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

	config.MustProcessDefault()

	if *server {
		db, err := dburl.Open(config.DefaultVars.XODBURL())
		if err != nil {
			return err
		}
		svc := service.New(db)
		server := http_server.New(db, svc)
		fmt.Printf("Listening on port %v\n", config.DefaultVars.HTTPPort)
		return http.ListenAndServe(
			fmt.Sprintf(":%d", config.DefaultVars.HTTPPort),
			server.NewRouter(),
		)
	}

	client, err := fetcher.NewClientFromURL("https://amaten.com")
	if err != nil {
		return err
	}
	giftCards, err := client.Fetch(context.Background(), &fetcher.FetchOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("discountRate = %v\n", giftCards[0].DiscountRate)
	fmt.Printf("salesPrice = %v\n", giftCards[0].SalesPrice)

	return nil
}
