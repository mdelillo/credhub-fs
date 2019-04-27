package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdelillo/credhub-fs/test/fake-uaa/handler"
	"github.com/mdelillo/credhub-fs/test/server"
)

type stringsFlag []string

func (i *stringsFlag) String() string {
	return fmt.Sprint(*i)
}
func (i *stringsFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	listenAddr    string
	certPath      string
	keyPath       string
	jwtSigningKey string
	clients       stringsFlag
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", "127.0.0.1:58443", "address to listen on")
	flag.StringVar(&certPath, "cert-path", "1270.0.1:58844", "path to TLS certificate")
	flag.StringVar(&keyPath, "key-path", "1270.0.1:58844", "path to TLS private key")
	flag.StringVar(&jwtSigningKey, "jwt-signing-key", "", "RSA key used to sign JWT tokens")
	flag.Var(&clients, "client", "client ID and secret, colon-separated, to be allowed for authentication. Can be specified multiple times.")
	flag.Parse()

	s := server.NewServer(
		listenAddr,
		certPath,
		keyPath,
		handler.NewUAAHandler(listenAddr, jwtSigningKey, clients),
	)

	go func() {
		if err := s.Start(); err != nil {
			log.Fatalf("Failed to start server: %s\n", err)
		}
	}()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	s.Shutdown()
}
