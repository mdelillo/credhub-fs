package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler"
	"github.com/mdelillo/credhub-fs/test/server"
)

var (
	listenAddr         string
	certPath           string
	keyPath            string
	authServerAddr     string
	jwtVerificationKey string
)

func main() {
	flag.StringVar(&listenAddr, "listen-addr", "127.0.0.1:58844", "address to listen on")
	flag.StringVar(&certPath, "cert-path", "1270.0.1:58844", "path to TLS certificate")
	flag.StringVar(&keyPath, "key-path", "1270.0.1:58844", "path to TLS private key")
	flag.StringVar(&authServerAddr, "auth-server-addr", "127.0.0.1:58443", "address of auth server")
	flag.StringVar(&jwtVerificationKey, "jwt-verification-key", "", "key used to verify JWT auth tokens")
	flag.Parse()

	s := server.NewServer(
		listenAddr,
		certPath,
		keyPath,
		handler.NewCredhubHandler(authServerAddr, jwtVerificationKey),
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
