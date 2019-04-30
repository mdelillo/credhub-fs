package credhub_test

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
	"github.com/mdelillo/credhub-fs/test/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		credhubServer           *ghttp.Server
		uaaServer               *ghttp.Server
		skipTLSVerifyHttpClient = &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		}
		emptyUUID              = uuid.MustParse("00000000-0000-0000-0000-000000000000")
		clientID               string
		clientSecret           string
		configureTokenHandlers func() string
	)

	BeforeEach(func() {
		credhubServer = ghttp.NewTLSServer()
		uaaServer = ghttp.NewTLSServer()
		clientID = helpers.RandomString()
		clientSecret = helpers.RandomString()

		configureTokenHandlers = func() string {
			token := helpers.RandomString()
			credhubServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/info"),
					ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"auth-server": {"url": "%s"}}`, uaaServer.URL())),
				),
			)
			uaaServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/oauth/token"),
					ghttp.VerifyFormKV("client_id", clientID),
					ghttp.VerifyFormKV("client_secret", clientSecret),
					ghttp.VerifyFormKV("grant_type", "client_credentials"),
					ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"access_token": "%s"}`, token)),
				),
			)
			return token
		}
	})

	AfterEach(func() {
		credhubServer.Close()
		uaaServer.Close()
	})

	Describe("GetCredentialByName", func() {
		It("returns the named credential", func() {
			credentialID := uuid.New()
			credentialName := "some-name"
			credentialType := "some-type"
			credentialValue := "some-value"
			credentialVersionCreatedAt := time.Now().UTC()

			token := configureTokenHandlers()

			credhubServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/data", "name="+credentialName),
					ghttp.VerifyHeaderKV("Authorization", "Bearer "+token),
					ghttp.RespondWith(http.StatusOK, fmt.Sprintf(
						`{"data": [{
							"id": "%s",
							"name": "%s",
							"type": "%s",
							"value": "%s",
							"version_created_at": "%s"
						}]}`, credentialID, credentialName, credentialType, credentialValue, credentialVersionCreatedAt.Format(time.RFC3339),
					)),
				),
			)

			credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
			client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
			credential, err := client.GetCredentialByName(credentialName)

			Expect(err).NotTo(HaveOccurred())
			Expect(credential.Name).To(Equal(credentialName))
			Expect(credential.ID).To(Equal(credentialID))
			Expect(credential.Type).To(Equal(credentialType))
			Expect(credential.Value).To(Equal(credentialValue))
			Expect(credential.VersionCreatedAt).To(BeTemporally("~", credentialVersionCreatedAt, time.Second))
		})

		Context("when getting the UAA URL fails", func() {
			It("returns an error", func() {
				client := credhub.NewClient("some-bad-url", clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName("some-name")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the credhub /info response is not 200", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusInternalServerError, `some-error`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName("some-name")
				Expect(err).To(MatchError(ContainSubstring(http.StatusText(http.StatusInternalServerError))))
			})
		})

		Context("when the credhub /info response is not valid JSON", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, `some-non-json-response`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName("some-name")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the UAA token request fails", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, `{"auth-server": {"url": "some-bad-url"}}`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName("some-name")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the UAA token response is not 200", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"auth-server": {"url": "%s"}}`, uaaServer.URL())),
					),
				)
				uaaServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/oauth/token"),
						ghttp.RespondWith(http.StatusInternalServerError, "some-error"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName("some-name")
				Expect(err).To(MatchError(ContainSubstring(http.StatusText(http.StatusInternalServerError))))
			})
		})

		Context("when the UAA token response is not valid JSON", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"auth-server": {"url": "%s"}}`, uaaServer.URL())),
					),
				)
				uaaServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/oauth/token"),
						ghttp.RespondWith(http.StatusOK, "some-non-json-response"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName("some-name")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the data response is 404", func() {
			It("returns an ErrCredentialNotFound", func() {
				credentialName := "some-name"

				configureTokenHandlers()

				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v1/data", "name="+credentialName),
						ghttp.RespondWith(http.StatusNotFound, "some-error"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName(credentialName)
				Expect(err).To(BeAssignableToTypeOf(&credhub.ErrCredentialNotFound{}))
			})
		})

		Context("when the data response is not 200 or 404", func() {
			It("returns an error", func() {
				credentialName := "some-name"

				configureTokenHandlers()

				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v1/data", "name="+credentialName),
						ghttp.RespondWith(http.StatusInternalServerError, "some-error"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName(credentialName)
				Expect(err).To(MatchError(ContainSubstring(http.StatusText(http.StatusInternalServerError))))
			})
		})

		Context("when the data response is not valid JSON", func() {
			It("returns an error", func() {
				credentialName := "some-name"

				configureTokenHandlers()

				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v1/data", "name="+credentialName),
						ghttp.RespondWith(http.StatusOK, "some-non-json-response"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName(credentialName)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the data response contains no credentials", func() {
			It("returns an error", func() {
				credentialName := "some-name"

				configureTokenHandlers()

				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v1/data", "name="+credentialName),
						ghttp.RespondWith(http.StatusOK, `{"data": []}`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName(credentialName)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the data response contains more than one credential", func() {
			It("returns an error", func() {
				credentialName := "some-name"

				configureTokenHandlers()

				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v1/data", "name="+credentialName),
						ghttp.RespondWith(http.StatusOK, `{"data": [{"name": "foo"}, {"name": "bar"}]}`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.GetCredentialByName(credentialName)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("FindCredentialsByPath", func() {
		It("returns credential names and versionCreatesAt dates at the given path", func() {
			path := "/some-path"
			credential1Name := path + "/some-name"
			credential1VersionCreatedAt := time.Now().UTC()
			credential2Name := path + "/some-other-name"
			credential2VersionCreatedAt := time.Now().UTC()

			token := configureTokenHandlers()

			credhubServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/data", "path="+path),
					ghttp.VerifyHeaderKV("Authorization", "Bearer "+token),
					ghttp.RespondWith(http.StatusOK, fmt.Sprintf(
						`{"credentials": [
							 {"name": "%s", "version_created_at": "%s"},
							 {"name": "%s", "version_created_at": "%s"}
						]}`,
						credential1Name, credential1VersionCreatedAt.Format(time.RFC3339),
						credential2Name, credential2VersionCreatedAt.Format(time.RFC3339),
					)),
				),
			)

			credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
			client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
			credentials, err := client.FindCredentialsByPath(path)

			Expect(err).NotTo(HaveOccurred())
			Expect(credentials).To(HaveLen(2))
			Expect(credentials[0].Name).To(Equal(credential1Name))
			Expect(credentials[0].VersionCreatedAt).To(BeTemporally("~", credential1VersionCreatedAt, time.Second))
			Expect(credentials[0].ID).To(Equal(emptyUUID))
			Expect(credentials[0].Type).To(BeEmpty())
			Expect(credentials[0].Value).To(BeEmpty())
			Expect(credentials[1].Name).To(Equal(credential2Name))
			Expect(credentials[1].VersionCreatedAt).To(BeTemporally("~", credential2VersionCreatedAt, time.Second))
			Expect(credentials[1].ID).To(Equal(emptyUUID))
			Expect(credentials[1].Type).To(BeEmpty())
			Expect(credentials[1].Value).To(BeEmpty())
		})

		Context("when getting the UAA URL fails", func() {
			It("returns an error", func() {
				client := credhub.NewClient("some-bad-url", clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath("some-path")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the credhub /info response is not 200", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusInternalServerError, `some-error`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath("some-path")
				Expect(err).To(MatchError(ContainSubstring(http.StatusText(http.StatusInternalServerError))))
			})
		})

		Context("when the credhub /info response is not valid JSON", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, `some-non-json-response`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath("some-path")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the UAA token request fails", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, `{"auth-server": {"url": "some-bad-url"}}`),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath("some-path")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the UAA token response is not 200", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"auth-server": {"url": "%s"}}`, uaaServer.URL())),
					),
				)
				uaaServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/oauth/token"),
						ghttp.RespondWith(http.StatusInternalServerError, "some-error"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath("some-path")
				Expect(err).To(MatchError(ContainSubstring(http.StatusText(http.StatusInternalServerError))))
			})
		})

		Context("when the UAA token response is not valid JSON", func() {
			It("returns an error", func() {
				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/info"),
						ghttp.RespondWith(http.StatusOK, fmt.Sprintf(`{"auth-server": {"url": "%s"}}`, uaaServer.URL())),
					),
				)
				uaaServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("POST", "/oauth/token"),
						ghttp.RespondWith(http.StatusOK, "some-non-json-response"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath("some-path")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the data response is not 200 or 404", func() {
			It("returns an error", func() {
				path := "some-path"

				configureTokenHandlers()

				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v1/data", "path="+path),
						ghttp.RespondWith(http.StatusNotFound, "some-error"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath(path)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the data response is not valid JSON", func() {
			It("returns an error", func() {
				path := "some-path"

				configureTokenHandlers()

				credhubServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v1/data", "path="+path),
						ghttp.RespondWith(http.StatusOK, "some-non-json-response"),
					),
				)

				credhubURL := strings.TrimPrefix(credhubServer.URL(), "https://")
				client := credhub.NewClient(credhubURL, clientID, clientSecret, skipTLSVerifyHttpClient)
				_, err := client.FindCredentialsByPath(path)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
