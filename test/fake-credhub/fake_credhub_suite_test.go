package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFakeCredhub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FakeCredhub Suite")
}
