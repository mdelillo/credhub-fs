package handler_test

import (
	"github.com/gin-gonic/gin"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/handler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Handler", func() {
	It("returns a handler that handles endpoints", func() {
		authServerURL := "some-auth-server-url"
		credhubHandler, err := handler.NewCredhubHandler(authServerURL, nil, nil)
		Expect(err).NotTo(HaveOccurred())

		routes := credhubHandler.(*gin.Engine).Routes()
		Expect(routes).To(ConsistOf(
			MatchFields(IgnoreExtras, Fields{"Path": Equal("/info"), "Method": Equal("GET")}),
			MatchFields(IgnoreExtras, Fields{"Path": Equal("/api/v1/data"), "Method": Equal("GET")}),
			MatchFields(IgnoreExtras, Fields{"Path": Equal("/api/v1/data"), "Method": Equal("PUT")}),
			MatchFields(IgnoreExtras, Fields{"Path": Equal("/api/v1/data"), "Method": Equal("DELETE")}),
		))
	})
})
