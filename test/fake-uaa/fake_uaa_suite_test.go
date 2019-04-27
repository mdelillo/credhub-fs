package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFakeUAA(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FakeUAA Suite")
}
