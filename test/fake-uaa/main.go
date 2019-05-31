package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/mdelillo/credhub-fs/test/fake-uaa/handler"
	"github.com/mdelillo/credhub-fs/test/server"
)

func main() {
	var opts struct {
		ListenAddr    string   `short:"l" long:"listen-addr" description:"address to listen on" default:"127.0.0.1:58443" required:"true"`
		CertPath      string   `short:"c" long:"cert-path" description:"path to TLS certificate" required:"true"`
		KeyPath       string   `short:"k" long:"key-path" description:"path to TLS private key" required:"true"`
		JWTSigningKey string   `short:"j" long:"jwt-signing-key" description:"RSA key used to sign JWT tokens" required:"true"`
		Clients       []string `long:"client" description:"client ID and secret, colon-separated, to be allowed for authentication. Can be specified multiple times."`
	}
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}

	uaaHandler, err := handler.NewUAAHandler(opts.ListenAddr, opts.JWTSigningKey, opts.Clients)
	if err != nil {
		log.Fatalf("Failed to create handler: %s\n", err.Error())
	}

	s := server.NewServer(opts.ListenAddr, opts.CertPath, opts.KeyPath, uaaHandler)

	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Failed to start server: %s\n", err.Error())
		}
	}()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	s.Shutdown()
}
