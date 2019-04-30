package ls_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ls Suite")
}
