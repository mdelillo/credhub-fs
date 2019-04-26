package credhub_fs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCredhubFs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CredhubFs Suite")
}
