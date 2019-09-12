package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/oinume/amamonitor/backend/cli"
	"github.com/oinume/amamonitor/backend/http_server"
)

func main() {
	m := &serverMain{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}
	if err := m.run(os.Args); err != nil {
		cli.WriteError(m.errStream, err)
		os.Exit(cli.ExitError)
	}
	os.Exit(cli.ExitOK)
}

type serverMain struct {
	outStream io.Writer
	errStream io.Writer
}

func (m *serverMain) run(args []string) error {
	flagSet := flag.NewFlagSet("server", flag.ContinueOnError)
	flagSet.SetOutput(m.errStream)
	var (
	//server = flagSet.Bool("server", false, "Run as a HTTP server")
	//logLevel        = flag.String("log-level", "info", "Log level")
	)
	if err := flagSet.Parse(args[1:]); err != nil {
		return err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}
	server := http_server.New()
	fmt.Printf("Listening on port %v\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), server.NewRouter())
}
