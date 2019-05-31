package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	"github.com/mdelillo/credhub-fs/test/fake-uaa/handler"
	"github.com/mdelillo/credhub-fs/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TokenHandler", func() {
	var (
		listenAddr    = "some-listen-addr"
		clientID      = "some-client-id"
		clientSecret  = "some-client-secret"
		clients       = []string{clientID + ":" + clientSecret}
		rsaKey        *rsa.PrivateKey
		jwtSigningKey string
	)

	BeforeEach(func() {
		var err error
		rsaKey, err = rsa.GenerateKey(rand.Reader, 4096)
		Expect(err).NotTo(HaveOccurred())
		jwtSigningKey = helpers.PrivateKeyToPEM(rsaKey)
	})

	It("responds with 200 and a valid token", func() {
		responseRecorder := httptest.NewRecorder()
		request := generateRequest(clientID, clientSecret, "client_credentials")

		uaaHandler, err := handler.NewUAAHandler(listenAddr, jwtSigningKey, clients)
		Expect(err).NotTo(HaveOccurred())
		uaaHandler.ServeHTTP(responseRecorder, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusOK))

		token := getRSATokenFromResponse(responseRecorder, rsaKey)
		Expect(token.Valid).To(BeTrue())
		Expect(token.Header["kid"]).To(Equal("legacy-token-key"))

		claims, ok := token.Claims.(jwt.MapClaims)
		Expect(ok).To(BeTrue())
		Expect(claims["client_id"]).To(Equal(clientID))
		Expect(claims["grant_type"]).To(Equal("client_credentials"))
		Expect(claims["iss"]).To(Equal(fmt.Sprintf("https://%s/oauth/token", listenAddr)))
		Expect(claims["scope"]).To(ConsistOf("credhub.read", "credhub.write"))
	})

	Context("when the grant type is not 'client_credentials'", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request := generateRequest(clientID, clientSecret, "some-invalid-grant-type")

			uaaHandler, err := handler.NewUAAHandler(listenAddr, jwtSigningKey, clients)
			Expect(err).NotTo(HaveOccurred())
			uaaHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "Grant type must be 'client_credentials'"}`))
		})
	})

	Context("when the 'client_id' is empty", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request := generateRequest("", clientSecret, "client_credentials")

			uaaHandler, err := handler.NewUAAHandler(listenAddr, jwtSigningKey, clients)
			Expect(err).NotTo(HaveOccurred())
			uaaHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "'client_id' must not be empty"}`))
		})
	})

	Context("when the 'client_secret' is empty", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request := generateRequest(clientID, "", "client_credentials")

			uaaHandler, err := handler.NewUAAHandler(listenAddr, jwtSigningKey, clients)
			Expect(err).NotTo(HaveOccurred())
			uaaHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "'client_secret' must not be empty"}`))
		})
	})

	Context("when the 'client_id' is not in the list of clients", func() {
		It("responds with a 401", func() {
			responseRecorder := httptest.NewRecorder()
			request := generateRequest("some-invalid-client-id", clientSecret, "client_credentials")

			uaaHandler, err := handler.NewUAAHandler(listenAddr, jwtSigningKey, clients)
			Expect(err).NotTo(HaveOccurred())
			uaaHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "incorrect 'client_id' and/or 'client_secret'"}`))
		})
	})

	Context("when the 'client_secret' does not match the 'client_id'", func() {
		It("responds with a 401", func() {
			responseRecorder := httptest.NewRecorder()
			request := generateRequest(clientID, "some-invalid-client-secret", "client_credentials")

			uaaHandler, err := handler.NewUAAHandler(listenAddr, jwtSigningKey, clients)
			Expect(err).NotTo(HaveOccurred())
			uaaHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "incorrect 'client_id' and/or 'client_secret'"}`))
		})
	})
})

func generateRequest(clientID, clientSecret, grantType string) *http.Request {
	request, err := http.NewRequest("POST", "/oauth/token", nil)
	request.PostForm = url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {grantType},
	}
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return request
}

func getRSATokenFromResponse(response *httptest.ResponseRecorder, rsaKey *rsa.PrivateKey) *jwt.Token {
	body, err := ioutil.ReadAll(response.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	ExpectWithOffset(1, json.Unmarshal(body, &tokenResponse)).To(Succeed())

	token, err := jwt.Parse(tokenResponse.AccessToken, func(token *jwt.Token) (interface{}, error) {
		_, signedUsingRSA := token.Method.(*jwt.SigningMethodRSA)
		ExpectWithOffset(1, signedUsingRSA).To(BeTrue())

		return &rsaKey.PublicKey, nil
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return token
}

func readBody(response *httptest.ResponseRecorder) []byte {
	body, err := ioutil.ReadAll(response.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return body
}
