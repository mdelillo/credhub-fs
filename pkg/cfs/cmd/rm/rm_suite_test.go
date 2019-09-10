package rm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rm Suite")
}
