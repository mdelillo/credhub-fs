package token_test

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/token"
	"github.com/mdelillo/credhub-fs/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {
	It("validates a token's signature and issuer", func() {
		iss := "some-iss"

		jwtSigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
		Expect(err).NotTo(HaveOccurred())
		keyString := helpers.PublicKeyToPEM(&jwtSigningKey.PublicKey)

		tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"iss": iss}).SignedString(jwtSigningKey)
		Expect(err).NotTo(HaveOccurred())

		claims := map[string]string{"iss": iss}
		validator, err := token.NewValidator(keyString)
		Expect(err).NotTo(HaveOccurred())
		Expect(validator.ValidateTokenWithClaims(tokenString, claims)).To(Succeed())
	})

	Context("when the signing key is invalid", func() {
		It("returns an error when creating the validator", func() {
			_, err := token.NewValidator("some-invalid-key")
			Expect(err).To(MatchError(ContainSubstring("failed to parse key:")))
		})
	})

	Context("when the token is not signed using RSA", func() {
		It("returns an error", func() {
			jwtSigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
			Expect(err).NotTo(HaveOccurred())
			keyString := helpers.PublicKeyToPEM(&jwtSigningKey.PublicKey)

			tokenString, err := jwt.New(jwt.SigningMethodHS256).SignedString([]byte("some-key"))
			Expect(err).NotTo(HaveOccurred())

			validator, err := token.NewValidator(keyString)
			Expect(err).NotTo(HaveOccurred())

			err = validator.ValidateTokenWithClaims(tokenString, nil)
			Expect(err).To(MatchError(ContainSubstring("token not signed using RSA")))
		})
	})

	Context("when the token cannot be parsed", func() {
		It("returns an error", func() {
			jwtSigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
			Expect(err).NotTo(HaveOccurred())
			keyString := helpers.PublicKeyToPEM(&jwtSigningKey.PublicKey)

			validator, err := token.NewValidator(keyString)
			Expect(err).NotTo(HaveOccurred())

			err = validator.ValidateTokenWithClaims("some-invalid-token", nil)
			Expect(err).To(MatchError(ContainSubstring("failed to parse token")))
		})
	})

	Context("when the token is signed by a different key", func() {
		It("returns an error", func() {
			jwtSigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
			Expect(err).NotTo(HaveOccurred())

			tokenString, err := jwt.New(jwt.SigningMethodRS256).SignedString(jwtSigningKey)
			Expect(err).NotTo(HaveOccurred())

			wrongJWTSigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
			Expect(err).NotTo(HaveOccurred())
			wrongKeyString := helpers.PublicKeyToPEM(&wrongJWTSigningKey.PublicKey)

			validator, err := token.NewValidator(wrongKeyString)
			Expect(err).NotTo(HaveOccurred())

			err = validator.ValidateTokenWithClaims(tokenString, nil)
			Expect(err).To(MatchError(ContainSubstring("failed to parse token")))
		})
	})

	Context("when the 'iss' claim does not match", func() {
		It("returns an error", func() {
			jwtSigningKey, err := rsa.GenerateKey(rand.Reader, 4096)
			Expect(err).NotTo(HaveOccurred())
			keyString := helpers.PublicKeyToPEM(&jwtSigningKey.PublicKey)

			tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"iss": "some-iss"}).SignedString(jwtSigningKey)
			Expect(err).NotTo(HaveOccurred())

			claims := map[string]string{"iss": "some-other-iss"}
			validator, err := token.NewValidator(keyString)
			Expect(err).NotTo(HaveOccurred())

			err = validator.ValidateTokenWithClaims(tokenString, claims)
			Expect(err).To(MatchError("iss claim does not match"))
		})
	})
})
