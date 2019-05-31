package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/credentials"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler/handlerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDataHandler", func() {
	var (
		credhubHandler      http.Handler
		fakeCredentialStore *handlerfakes.FakeCredentialStore
		fakeTokenValidator  *handlerfakes.FakeTokenValidator
		authServerURL       = "some-auth-server-url"
	)

	BeforeEach(func() {
		fakeCredentialStore = &handlerfakes.FakeCredentialStore{}
		fakeTokenValidator = &handlerfakes.FakeTokenValidator{}

		var err error
		credhubHandler, err = handler.NewCredhubHandler(authServerURL, fakeCredentialStore, fakeTokenValidator)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when 'name' is provided", func() {
		It("gets a credential from the store by name", func() {
			name := "some-name"
			expectedCredential := credentials.Credential{
				ID:               uuid.New(),
				Name:             name,
				Type:             "some-type",
				Value:            "some-value",
				VersionCreatedAt: time.Now().UTC(),
			}
			fakeCredentialStore.GetByNameReturns(expectedCredential, true)

			responseRecorder := httptest.NewRecorder()
			request := getDataByNameRequest(name, "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))

			var response struct {
				Data []credentials.Credential
			}
			Expect(json.Unmarshal(readBody(responseRecorder), &response)).To(Succeed())
			Expect(response.Data).To(HaveLen(1))
			Expect(response.Data[0]).To(Equal(expectedCredential))

			Expect(fakeCredentialStore.GetByNameArgsForCall(0)).To(Equal(name))
		})

		Context("when the credential does not exist", func() {
			It("responds with a 404", func() {
				fakeCredentialStore.GetByNameReturns(credentials.Credential{}, false)

				responseRecorder := httptest.NewRecorder()
				request := getDataByNameRequest("some-nonexistent-name", "some-token")

				credhubHandler.ServeHTTP(responseRecorder, request)

				Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
				Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "The request could not be completed because the credential does not exist or you do not have sufficient authorization."}`))
			})
		})
	})

	Context("when 'path' is provided", func() {
		It("gets credentials from the store by path", func() {
			path := "some-path"
			expectedCredentials := []credentials.Credential{
				{
					ID:               uuid.New(),
					Name:             "some-credential",
					Type:             "some-type",
					Value:            "some-value",
					VersionCreatedAt: time.Now().UTC(),
				},
				{
					ID:               uuid.New(),
					Name:             "some-other-credential",
					Type:             "some-type",
					Value:            "some-value",
					VersionCreatedAt: time.Now().UTC(),
				},
			}
			fakeCredentialStore.GetByPathReturns(expectedCredentials)

			responseRecorder := httptest.NewRecorder()
			request := getDataByPathRequest(path, "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))

			Expect(readBody(responseRecorder)).To(MatchJSON(fmt.Sprintf(
				`{
					 "credentials": [
					   {"name": "%s", "version_created_at": "%s"},
					   {"name": "%s", "version_created_at": "%s"}
				   ]
				 }`,
				expectedCredentials[0].Name, expectedCredentials[0].VersionCreatedAt.Format(time.RFC3339Nano),
				expectedCredentials[1].Name, expectedCredentials[1].VersionCreatedAt.Format(time.RFC3339Nano),
			)))

			Expect(fakeCredentialStore.GetByPathArgsForCall(0)).To(Equal(path))
		})
	})

	Context("when both 'name' and 'path' are provided", func() {
		It("uses the 'path'", func() {
			path := "some-path"

			responseRecorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/api/v1/data?name=some-name&path="+path, nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Add("Authorization", "Bearer some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))

			Expect(fakeCredentialStore.GetByPathCallCount()).To(Equal(1))
			Expect(fakeCredentialStore.GetByPathArgsForCall(0)).To(Equal(path))
			Expect(fakeCredentialStore.GetByNameCallCount()).To(Equal(0))
		})
	})

	Context("when neither 'name' or 'path' are provided", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request, err := http.NewRequest("GET", "/api/v1/data", nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Add("Authorization", "Bearer some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "The query parameter name is required for this request."}`))
		})
	})

	Context("when no Authorization header is set", func() {
		It("responds with a 401", func() {
			responseRecorder := httptest.NewRecorder()
			request := getDataByNameRequest("some-name", "")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{
				"error": "invalid_token",
				"error_description": "Full authentication is required to access this resource"
			}`))
		})
	})

	Context("when token validation fails", func() {
		It("responds with a 401", func() {
			fakeTokenValidator.ValidateTokenWithClaimsReturns(errors.New("some-error"))

			responseRecorder := httptest.NewRecorder()
			request := getDataByNameRequest("some-name", "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{
				"error": "invalid_token",
				"error_description": "The request token is malformed. Please validate that your request token was issued by the UAA server authorized by CredHub."
			}`))
		})
	})
})

func getDataByNameRequest(name string, token string) *http.Request {
	request, err := http.NewRequest("GET", "/api/v1/data?name="+name, nil)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	return request
}

func getDataByPathRequest(path string, token string) *http.Request {
	request, err := http.NewRequest("GET", "/api/v1/data?path="+path, nil)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	return request
}
