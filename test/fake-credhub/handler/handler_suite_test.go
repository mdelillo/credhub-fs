package handler_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FakeCredhub Handler Suite")
}

var _ = BeforeSuite(func() {
	gin.DefaultWriter = GinkgoWriter
})

func readBody(response *httptest.ResponseRecorder) []byte {
	body, err := ioutil.ReadAll(response.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return body
}
