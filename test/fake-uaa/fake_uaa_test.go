package main_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mdelillo/credhub-fs/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("FakeUAA", func() {
	var (
		fakeUAA       string
		listenAddr    string
		jwtSigningKey *rsa.PrivateKey
		clientID      string
		clientSecret  string
	)

	BeforeSuite(func() {
		var err error
		fakeUAA, err = gexec.Build(filepath.Join("github.com", "mdelillo", "credhub-fs", "test", "fake-uaa"))
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		var err error
		jwtSigningKey, err = rsa.GenerateKey(rand.Reader, 4096)
		Expect(err).NotTo(HaveOccurred())

		listenAddr = helpers.GetFreeAddr()
		clientID = helpers.RandomString()
		clientSecret = helpers.RandomString()
		cmd := exec.Command(
			fakeUAA,
			"--listen-addr", listenAddr,
			"--cert-path", filepath.Join("..", "fixtures", "127.0.0.1-cert.pem"),
			"--key-path", filepath.Join("..", "fixtures", "127.0.0.1-key.pem"),
			"--jwt-signing-key", helpers.PrivateKeyToPEM(jwtSigningKey),
			"--client", clientID+":"+clientSecret,
		)
		_, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(helpers.WaitForServerToBeAvailable(listenAddr, 5*time.Second)).To(Succeed())
	})

	AfterEach(func() {
		gexec.KillAndWait()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	It("grants a JWT token signed by the RSA signing key", func() {
		tokenUrl := fmt.Sprintf("https://%s/oauth/token", listenAddr)
		values := url.Values{
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"grant_type":    {"client_credentials"},
		}
		resp, err := helpers.HTTPClient.PostForm(tokenUrl, values)
		Expect(err).NotTo(HaveOccurred())

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		var tokenResponse struct {
			AccessToken string `json:"access_token"`
		}
		Expect(json.Unmarshal(body, &tokenResponse)).To(Succeed())

		token, err := jwt.Parse(tokenResponse.AccessToken, func(token *jwt.Token) (interface{}, error) { return &jwtSigningKey.PublicKey, nil })
		Expect(err).NotTo(HaveOccurred())
		Expect(token.Valid).To(BeTrue())
	})
})
