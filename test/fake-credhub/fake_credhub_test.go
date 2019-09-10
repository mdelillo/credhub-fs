package main_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/mdelillo/credhub-fs/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("FakeCredhub", func() {
	var (
		fakeCredhub    string
		listenAddr     string
		authServerAddr = "some-auth-server-addr"
		jwtSigningKey  *rsa.PrivateKey
		makeRequest    func(method, path, body, authToken string) (statusCode int, responseBody string)
		get            func(path, authToken string) (statusCode int, responseBody string)
		put            func(path, body, authToken string) (statusCode int, responseBody string)
		delete         func(path, authToken string) (statusCode int, responseBody string)
	)

	BeforeSuite(func() {
		var err error
		fakeCredhub, err = gexec.Build(filepath.Join("github.com", "mdelillo", "credhub-fs", "test", "fake-credhub"))
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		var err error
		jwtSigningKey, err = rsa.GenerateKey(rand.Reader, 4096)
		Expect(err).NotTo(HaveOccurred())

		listenAddr = helpers.GetFreeAddr()
		cmd := exec.Command(
			fakeCredhub,
			"--listen-addr", listenAddr,
			"--cert-path", filepath.Join("..", "fixtures", "127.0.0.1-cert.pem"),
			"--key-path", filepath.Join("..", "fixtures", "127.0.0.1-key.pem"),
			"--auth-server-addr", authServerAddr,
			"--jwt-verification-key", helpers.PublicKeyToPEM(&jwtSigningKey.PublicKey),
		)
		_, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(helpers.WaitForServerToBeAvailable(listenAddr, 5*time.Second)).To(Succeed())

		makeRequest = func(method, path, body, authToken string) (int, string) {
			url := fmt.Sprintf("https://%s/%s", listenAddr, path)
			req, err := http.NewRequest(method, url, strings.NewReader(body))
			Expect(err).NotTo(HaveOccurred())

			req.Header.Add("Accept", "application/json")
			req.Header.Add("Content-Type", "application/json")
			if authToken != "" {
				req.Header.Add("Authorization", "Bearer "+authToken)
			}

			resp, err := helpers.HTTPClient.Do(req)
			Expect(err).NotTo(HaveOccurred())

			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			return resp.StatusCode, string(respBody)
		}
		get = func(path, authToken string) (int, string) {
			return makeRequest(http.MethodGet, path, "", authToken)
		}
		put = func(path, body, authToken string) (int, string) {
			return makeRequest(http.MethodPut, path, body, authToken)
		}
		delete = func(path, authToken string) (int, string) {
			return makeRequest(http.MethodDelete, path, "", authToken)
		}
	})

	AfterEach(func() {
		gexec.KillAndWait()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	It("implements GET /info", func() {
		statusCode, body := get("info", "")
		Expect(statusCode).To(Equal(http.StatusOK))

		Expect(string(body)).To(MatchJSON(fmt.Sprintf(
			`{
				"auth-server": {"url": "%s"},
				"app": {"name": "Fake CredHub"}
			}`,
			authServerAddr,
		)))
	})

	Describe("/api/v1/data", func() {
		It("can set and get value type credentials when authenticated", func() {
			name := "/" + helpers.RandomString()
			value := helpers.RandomString()

			setCredReqBody := fmt.Sprintf(`{"name": "%s", "value": "%s", "type": "value"}`, name, value)
			token := generateJWTToken(authServerAddr, jwtSigningKey)

			setCredStatusCode, setCredResponseBody := put("api/v1/data", setCredReqBody, token)
			Expect(setCredStatusCode).To(Equal(http.StatusOK))

			var credFromSet credential
			Expect(json.Unmarshal([]byte(setCredResponseBody), &credFromSet)).To(Succeed())
			Expect(credFromSet.ID).NotTo(BeEmpty())
			Expect(credFromSet.VersionCreatedAt).To(BeTemporally("~", time.Now(), 5*time.Second))
			Expect(credFromSet.Name).To(Equal(name))
			Expect(credFromSet.Value).To(Equal(value))
			Expect(credFromSet.Type).To(Equal("value"))

			getCredStatusCode, getCredResponseBody := get("api/v1/data?name="+name, token)
			Expect(getCredStatusCode).To(Equal(http.StatusOK))

			var resp getCredResponse
			Expect(json.Unmarshal([]byte(getCredResponseBody), &resp)).To(Succeed())
			Expect(resp.Data).To(HaveLen(1))

			credFromGet := resp.Data[0]
			Expect(credFromGet).To(Equal(credFromSet))
		})

		It("can list credential names and versionCreatedAt dates at a specific path", func() {
			topLevelName := "/" + helpers.RandomString()
			nestedName := "/some-dir/" + helpers.RandomString()
			anotherNestedName := "/some-dir/" + helpers.RandomString()
			deeplyNestedName := "/some-dir/some-nested-dir/" + helpers.RandomString()

			token := generateJWTToken(authServerAddr, jwtSigningKey)
			for _, name := range []string{topLevelName, nestedName, anotherNestedName, deeplyNestedName} {
				body := fmt.Sprintf(`{"name": "%s", "value": "some-value", "type": "value"}`, name)
				statusCode, _ := put("api/v1/data", body, token)
				Expect(statusCode).To(Equal(http.StatusOK))
			}

			By("listing a path with credentials in it")
			statusCode, respBody := get("api/v1/data?path=/some-dir", token)
			Expect(statusCode).To(Equal(http.StatusOK))

			var credentials listCredsResponse
			Expect(json.Unmarshal([]byte(respBody), &credentials)).To(Succeed())

			var credentialNames []string
			for _, credential := range credentials.Credentials {
				credentialNames = append(credentialNames, credential.Name)
				Expect(credential.VersionCreatedAt).To(BeTemporally("~", time.Now(), 5*time.Second))
			}
			Expect(credentialNames).To(ConsistOf(nestedName, anotherNestedName, deeplyNestedName))

			By("listing a credential directly")
			statusCode, respBody = get("api/v1/data?path="+topLevelName, token)
			Expect(statusCode).To(Equal(http.StatusOK))
			Expect(statusCode).To(Equal(http.StatusOK))
			Expect(json.Unmarshal([]byte(respBody), &credentials)).To(Succeed())
			Expect(credentials.Credentials).To(BeEmpty())
		})

		It("can delete credentials", func() {
			name1 := "/" + helpers.RandomString()
			name2 := "/" + helpers.RandomString()

			token := generateJWTToken(authServerAddr, jwtSigningKey)
			for _, name := range []string{name1, name2} {
				body := fmt.Sprintf(`{"name": "%s", "value": "some-value", "type": "value"}`, name)
				statusCode, _ := put("api/v1/data", body, token)
				Expect(statusCode).To(Equal(http.StatusOK))
			}

			var credentials listCredsResponse
			statusCode, respBody := get("api/v1/data?path=/", token)
			Expect(statusCode).To(Equal(http.StatusOK))
			Expect(json.Unmarshal([]byte(respBody), &credentials)).To(Succeed())
			Expect(credentials.Credentials).To(HaveLen(2))

			statusCode, respBody = delete("api/v1/data?name="+name1, token)
			Expect(statusCode).To(Equal(http.StatusNoContent))
			Expect(respBody).To(BeEmpty())

			statusCode, respBody = get("api/v1/data?path=/", token)
			Expect(statusCode).To(Equal(http.StatusOK))
			Expect(json.Unmarshal([]byte(respBody), &credentials)).To(Succeed())
			Expect(credentials.Credentials).To(HaveLen(1))

			statusCode, respBody = delete("api/v1/data?name="+name2, token)
			Expect(statusCode).To(Equal(http.StatusNoContent))
			Expect(respBody).To(BeEmpty())

			statusCode, respBody = get("api/v1/data?path=/", token)
			Expect(statusCode).To(Equal(http.StatusOK))
			Expect(json.Unmarshal([]byte(respBody), &credentials)).To(Succeed())
			Expect(credentials.Credentials).To(BeEmpty())
		})

		It("returns 401s when no token is provided", func() {
			name := "/" + helpers.RandomString()
			value := helpers.RandomString()
			setCredReqBody := fmt.Sprintf(`{"name": "%s", "value": "%s", "type": "value"}`, name, value)

			setCredStatusCode, setCredResponseBody := put("api/v1/data", setCredReqBody, "")
			Expect(setCredStatusCode).To(Equal(http.StatusUnauthorized))
			Expect(setCredResponseBody).To(MatchJSON(`{"error": "invalid_token", "error_description": "Full authentication is required to access this resource"}`))

			getCredStatusCode, getCredResponseBody := get("api/v1/data?name="+name, "")
			Expect(getCredStatusCode).To(Equal(http.StatusUnauthorized))
			Expect(getCredResponseBody).To(MatchJSON(`{"error": "invalid_token", "error_description": "Full authentication is required to access this resource"}`))
		})

		It("returns 401s when the token is issued by a different auth server", func() {
			name := "/" + helpers.RandomString()
			value := helpers.RandomString()
			setCredReqBody := fmt.Sprintf(`{"name": "%s", "value": "%s", "type": "value"}`, name, value)
			badToken := generateJWTToken("some-other-auth-server-addr", jwtSigningKey)

			setCredStatusCode, setCredResponseBody := put("api/v1/data", setCredReqBody, badToken)
			Expect(setCredStatusCode).To(Equal(http.StatusUnauthorized))
			Expect(setCredResponseBody).To(MatchJSON(`{"error": "invalid_token", "error_description": "The request token is malformed. Please validate that your request token was issued by the UAA server authorized by CredHub."}`))

			getCredStatusCode, getCredResponseBody := get("api/v1/data?name="+name, badToken)
			Expect(getCredStatusCode).To(Equal(http.StatusUnauthorized))
			Expect(getCredResponseBody).To(MatchJSON(`{"error": "invalid_token", "error_description": "The request token is malformed. Please validate that your request token was issued by the UAA server authorized by CredHub."}`))
		})
	})
})

type credential struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Value            string    `json:"value"`
	VersionCreatedAt time.Time `json:"version_created_at"`
}

type setCredResponse credential
type getCredResponse struct {
	Data []credential `json:"data"`
}
type listCredsResponse struct {
	Credentials []struct {
		Name             string    `json:"name"`
		VersionCreatedAt time.Time `json:"version_created_at"`
	} `json:"credentials"`
}

func generateJWTToken(authServerAddr string, signingKey *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"client_id":  "some-client-id",
		"grant_type": "client_credentials",
		"iss":        authServerAddr + "/oauth/token",
		"scope":      []string{"credhub.read", "credhub.write"},
	})
	token.Header["kid"] = "legacy-token-key"
	tokenString, err := token.SignedString(signingKey)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return tokenString
}
