package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler/handlerfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeleteDataHandler", func() {
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

	It("deletes a credential from the store by name", func() {
		name := "some-name"
		fakeCredentialStore.DeleteReturns(true)

		responseRecorder := httptest.NewRecorder()
		request := deleteRequest(name, "some-token")

		credhubHandler.ServeHTTP(responseRecorder, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusNoContent))
		Expect(fakeCredentialStore.DeleteCallCount()).To(Equal(1))
		Expect(fakeCredentialStore.DeleteArgsForCall(0)).To(Equal(name))
	})

	Context("when the credential does not exist", func() {
		It("responds with a 404", func() {
			fakeCredentialStore.DeleteReturns(false)

			responseRecorder := httptest.NewRecorder()
			request := deleteRequest("some-nonexistent-name", "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusNotFound))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{"error": "The request could not be completed because the credential does not exist or you do not have sufficient authorization."}`))
		})
	})

	Context("when a name is not provided", func() {
		It("responds with a 400", func() {
			responseRecorder := httptest.NewRecorder()
			request, err := http.NewRequest("DELETE", "/api/v1/data", nil)
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
			request := deleteRequest("some-name", "")

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
			request := deleteRequest("some-name", "some-token")

			credhubHandler.ServeHTTP(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnauthorized))
			Expect(readBody(responseRecorder)).To(MatchJSON(`{
				"error": "invalid_token",
				"error_description": "The request token is malformed. Please validate that your request token was issued by the UAA server authorized by CredHub."
			}`))
		})
	})
})

func deleteRequest(name string, token string) *http.Request {
	request, err := http.NewRequest("DELETE", "/api/v1/data?name="+name, nil)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	if token != "" {
		request.Header.Add("Authorization", "Bearer "+token)
	}

	return request
}
