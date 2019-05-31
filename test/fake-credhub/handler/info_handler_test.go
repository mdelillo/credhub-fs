package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InfoHandler", func() {
	It("responds with 200 and the auth-server URL and app name", func() {
		authServerURL := "some-auth-server-url"
		credhubHandler, err := handler.NewCredhubHandler(authServerURL, nil, nil)
		Expect(err).NotTo(HaveOccurred())

		responseRecorder := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "/info", nil)
		Expect(err).NotTo(HaveOccurred())

		credhubHandler.ServeHTTP(responseRecorder, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		Expect(readBody(responseRecorder)).To(MatchJSON(fmt.Sprintf(
			`{"auth-server": {"url": "%s"}, "app": {"name": "Fake CredHub"}}`, authServerURL,
		)))
	})
})
