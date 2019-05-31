package handler_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler/handlerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PutDataHandler", func() {
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

	It("sets a credential in the store", func() {
		name := "some-name"
		credType := "value"
		value := "some-value"

		responseRecorder := httptest.NewRecorder()
		request := setRequest(name, credType, value, "some-token")

		credhubHandler.ServeHTTP(responseRecorder, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusOK))

		cred := fakeCredentialStore.SetArgsForCall(0)
		Expect(cred.Name).To(Equal(name))
		Expect(cred.Type).To(Equal(credType))
		Expect(cred.Value).To(Equal(value))
		Expect(cred.ID).NotTo(Equal(uuid.Nil))
		Expect(cred.VersionCreatedAt).To(BeTemporally("~", time.Now().UTC(), 5*time.Second))
	})

	Context("when the request body is not JSON", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()

			request, err := http.NewRequest("PUT", "/api/v1/data", strings.NewReader("some-non-json-body"))
			Expect(err).NotTo(HaveOccurred())
			request.Header.Add("Authorization", "Bearer some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{
				"error": "The request could not be fulfilled because the request path or body did not meet expectation. Please check the documentation for required formatting and retry your request."
			}`))
		})
	})

	Context("when the type is not 'value'", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request := setRequest("some-name", "some-non-value-type", "some-value", "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{
				"error": "Only 'value' types are supported"
			}`))
		})
	})

	Context("when the name is empty", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request := setRequest("", "value", "some-value", "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when the value is empty", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request := setRequest("some-name", "value", "", "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when no Authorization header is set", func() {
		It("responds with a 401", func() {
			responseRecorder := httptest.NewRecorder()

			request := setRequest("some-name", "value", "some-value", "")

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
			request := setRequest("some-name", "value", "some-value", "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{
				"error": "invalid_token",
				"error_description": "The request token is malformed. Please validate that your request token was issued by the UAA server authorized by CredHub."
			}`))
		})
	})
})

func setRequest(name, credType, value, token string) *http.Request {
	body := strings.NewReader(fmt.Sprintf(`{"name": "%s", "type": "%s", "value": "%s"}`, name, credType, value))
	request, err := http.NewRequest("PUT", "/api/v1/data", body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	return request
}
