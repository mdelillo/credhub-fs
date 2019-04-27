package credhub_fs_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mdelillo/credhub-fs/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("CredhubFs", func() {
	var (
		cfs               string
		fakeCredhub       string
		credhubListenAddr string
		credhubCertPath   string
		credhubKeyPath    string
		fakeUAA           string
		uaaListenAddr     string
		uaaCertPath       string
		uaaKeyPath        string
		clientID          string
		clientSecret      string
		jwtSigningKey     *rsa.PrivateKey
		setValueInCredhub func(name, value string)
	)

	BeforeSuite(func() {
		var err error
		cfs, err = gexec.Build(filepath.Join("github.com", "mdelillo", "credhub-fs", "cmd", "cfs"))
		Expect(err).NotTo(HaveOccurred())

		fakeCredhub, err = gexec.Build(filepath.Join("github.com", "mdelillo", "credhub-fs", "test", "fake-credhub"))
		Expect(err).NotTo(HaveOccurred())

		fakeUAA, err = gexec.Build(filepath.Join("github.com", "mdelillo", "credhub-fs", "test", "fake-uaa"))
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		var err error
		jwtSigningKey, err = rsa.GenerateKey(rand.Reader, 4096)
		Expect(err).NotTo(HaveOccurred())

		credhubListenAddr = helpers.GetFreeAddr()
		credhubCertPath, credhubKeyPath = helpers.GenerateSelfSignedCert("127.0.0.1")
		uaaListenAddr = helpers.GetFreeAddr()
		uaaCertPath, uaaKeyPath = helpers.GenerateSelfSignedCert("127.0.0.1")

		cmd := exec.Command(
			fakeCredhub,
			"--listen-addr", credhubListenAddr,
			"--cert-path", credhubCertPath,
			"--key-path", credhubKeyPath,
			"--auth-server-addr", "https://"+uaaListenAddr,
			"--jwt-verification-key", helpers.PublicKeyToPEM(&jwtSigningKey.PublicKey),
		)
		_, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(helpers.WaitForServerToBeAvailable(credhubListenAddr, 5*time.Second)).To(Succeed())

		clientID = helpers.RandomString()
		clientSecret = helpers.RandomString()
		cmd = exec.Command(
			fakeUAA,
			"--listen-addr", uaaListenAddr,
			"--cert-path", uaaCertPath,
			"--key-path", uaaKeyPath,
			"--jwt-signing-key", helpers.PrivateKeyToPEM(jwtSigningKey),
			"--client", clientID+":"+clientSecret,
		)
		_, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(helpers.WaitForServerToBeAvailable(uaaListenAddr, 5*time.Second)).To(Succeed())

		getUAAToken := func() string {
			tokenURL := fmt.Sprintf("https://%s/oauth/token", uaaListenAddr)
			values := url.Values{
				"client_id":     {clientID},
				"client_secret": {clientSecret},
				"grant_type":    {"client_credentials"},
			}
			resp, err := helpers.HTTPClient.PostForm(tokenURL, values)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			var tokenResponse struct {
				AccessToken string `json:"access_token"`
			}
			ExpectWithOffset(1, json.Unmarshal(body, &tokenResponse)).To(Succeed())

			return tokenResponse.AccessToken
		}

		setValueInCredhub = func(name, value string) {
			url := fmt.Sprintf("https://%s/api/v1/data", credhubListenAddr)
			body := strings.NewReader(fmt.Sprintf(
				`{"name": "%s", "value": "%s", "type": "value"}`,
				name,
				value,
			))

			req, err := http.NewRequest(http.MethodPut, url, body)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			req.Header.Add("Authorization", "Bearer "+getUAAToken())

			resp, err := helpers.HTTPClient.Do(req)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			ExpectWithOffset(1, resp.StatusCode).To(Equal(http.StatusOK), string(respBody))
		}
	})

	AfterEach(func() {
		gexec.KillAndWait()
		os.Remove(credhubCertPath)
		os.Remove(credhubKeyPath)
		os.Remove(uaaCertPath)
		os.Remove(uaaKeyPath)
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	It("prints the help text", func() {
		cmd := exec.Command(cfs)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		Expect(session).To(gbytes.Say("cfs interacts with CredHub using Unix filesystem commands"))
	})

	Describe("cfs cat", func() {
		It("shows the value of a credential", func() {
			name := "/" + helpers.RandomString()
			value := helpers.RandomString()
			setValueInCredhub(fmt.Sprintf("%s", name), value)

			cmd := exec.Command(cfs, "cat", name)
			cmd.Env = append(cmd.Env, "CREDHUB_ADDR="+credhubListenAddr)
			cmd.Env = append(cmd.Env, "CLIENT_ID="+clientID)
			cmd.Env = append(cmd.Env, "CLIENT_SECRET="+clientSecret)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))

			Expect(session).To(gbytes.Say(value))
		})
	})
})
