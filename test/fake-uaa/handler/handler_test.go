package handler_test

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/gin-gonic/gin"
	"github.com/mdelillo/credhub-fs/test/fake-uaa/handler"
	"github.com/mdelillo/credhub-fs/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	It("returns a handler that handles POST /oauth/token", func() {
		listenAddr := "some-listen-addr"
		clients := []string{"some-client-id:some-secret"}

		rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
		Expect(err).NotTo(HaveOccurred())
		jwtSigningKey := helpers.PrivateKeyToPEM(rsaKey)

		uaaHandler, err := handler.NewUAAHandler(listenAddr, jwtSigningKey, clients)
		Expect(err).NotTo(HaveOccurred())

		routes := uaaHandler.(*gin.Engine).Routes()
		Expect(routes).To(HaveLen(1))
		Expect(routes[0].Path).To(Equal("/oauth/token"))
		Expect(routes[0].Method).To(Equal("POST"))
	})

	Context("when the clients list is not formatted properly", func() {
		It("returns an error", func() {
			listenAddr := "some-listen-addr"
			invalidClients := []string{"valid:valid", "invalid"}

			rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
			Expect(err).NotTo(HaveOccurred())
			jwtSigningKey := helpers.PrivateKeyToPEM(rsaKey)

			_, err = handler.NewUAAHandler(listenAddr, jwtSigningKey, invalidClients)
			Expect(err).To(MatchError("'clients' must contain colon-separated client IDs and secrets"))
		})
	})

	Context("when the JWT signing key is invalid", func() {
		It("returns an error", func() {
			listenAddr := "some-listen-addr"
			clients := []string{"some-client-id:some-secret"}
			invalidJWTSigningKey := "some-invalid-jwt-signing-key"

			_, err := handler.NewUAAHandler(listenAddr, invalidJWTSigningKey, clients)
			Expect(err).To(MatchError(ContainSubstring("failed to parse JWT signing key: ")))
		})
	})
})
