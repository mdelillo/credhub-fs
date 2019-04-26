package credhub_fs_test

import (
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("CredhubFs", func() {
	var cfs string

	BeforeSuite(func() {
		var err error
		cfs, err = gexec.Build(filepath.Join("github.com", "mdelillo", "credhub-fs", "cmd", "cfs"))
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		gexec.KillAndWait()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	It("prints the help text", func() {
		cmd := exec.Command(cfs)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))
		Expect(session).To(gbytes.Say("cfs interacts with CredHub using Unix filesystem commands"))
	})
})
