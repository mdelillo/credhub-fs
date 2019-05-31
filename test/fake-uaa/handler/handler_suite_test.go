package handler_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FakeUAA Handler Suite")
}

var _ = BeforeSuite(func() {
	gin.DefaultWriter = GinkgoWriter
})
