package cat_test

import (
	"bytes"
	"errors"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/cat"
	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/cat/catfakes"
	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
)

var _ = Describe("Cat", func() {
	var fakeCredhubClient *catfakes.FakeCredhubClient
	var dependencies cmdutil.Dependencies

	BeforeEach(func() {
		fakeCredhubClient = &catfakes.FakeCredhubClient{}
		dependencies = cmdutil.NewDependencies()
		dependencies.SetCredhubClient(fakeCredhubClient)
	})

	It("prints the value of a credential", func() {
		value := "some-value"
		fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{Value: value}, nil)

		path := "/some/path/to/cred"
		var output bytes.Buffer
		cmd := cat.NewCmdCat(dependencies)
		cmd.SetArgs([]string{path})
		cmd.SetOutput(&output)

		Expect(cmd.Execute()).To(Succeed())

		Expect(output.String()).To(Equal(value + "\n"))
		Expect(fakeCredhubClient.GetCredentialByNameArgsForCall(0)).To(Equal(path))
	})

	Context("when no credential is found", func() {
		It("returns an error", func() {
			fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, &credhub.ErrCredentialNotFound{})

			var output bytes.Buffer
			cmd := cat.NewCmdCat(dependencies)
			cmd.SetArgs([]string{"/some-path"})
			cmd.SetOutput(&output)

			Expect(cmd.Execute()).To(MatchError("'/some-path': no such credential or path"))
			Expect(cmd.SilenceUsage).To(BeTrue())
		})
	})

	Context("when no arguments are provided", func() {
		It("returns an error and shows the usage", func() {
			cmd := cat.NewCmdCat(dependencies)
			cmd.SetArgs([]string{})
			cmd.SetOutput(ioutil.Discard)

			Expect(cmd.Execute()).To(MatchError(ContainSubstring("must provide a credential path")))
			Expect(cmd.SilenceUsage).To(BeFalse())
		})
	})

	Context("when getting the credential fails", func() {
		It("returns an error", func() {
			fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, errors.New("some-error"))

			cmd := cat.NewCmdCat(dependencies)
			cmd.SetArgs([]string{"some-cred"})
			cmd.SetOutput(ioutil.Discard)

			Expect(cmd.Execute()).To(MatchError(ContainSubstring("failed to get credential: some-error")))
			Expect(cmd.SilenceUsage).To(BeTrue())
		})
	})
})
