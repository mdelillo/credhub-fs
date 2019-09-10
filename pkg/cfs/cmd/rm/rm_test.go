package rm_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/rm"
	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/rm/rmfakes"
	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
)

var _ = Describe("Rm", func() {
	var fakeCredhubClient *rmfakes.FakeCredhubClient
	var dependencies cmdutil.Dependencies

	BeforeEach(func() {
		fakeCredhubClient = &rmfakes.FakeCredhubClient{}
		dependencies = cmdutil.NewDependencies()
		dependencies.SetCredhubClient(fakeCredhubClient)
	})

	Context("when the path matches a single credential", func() {
		It("removes the credential", func() {
			path := "/some/path/to/cred"
			fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{Name: path}, nil)

			var output bytes.Buffer
			cmd := rm.NewCmdRm(dependencies)
			cmd.SetArgs([]string{path})
			cmd.SetOutput(&output)

			Expect(cmd.Execute()).To(Succeed())

			Expect(fakeCredhubClient.GetCredentialByNameCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.GetCredentialByNameArgsForCall(0)).To(Equal(path))
			Expect(fakeCredhubClient.DeleteCredentialByNameCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.DeleteCredentialByNameArgsForCall(0)).To(Equal(path))
		})

		Context("when deleting the credential fails", func() {
			It("returns an error", func() {
				path := "/some/path/to/cred"
				fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{Name: path}, nil)
				fakeCredhubClient.DeleteCredentialByNameReturns(errors.New("some-error"))

				var output bytes.Buffer
				cmd := rm.NewCmdRm(dependencies)
				cmd.SetArgs([]string{path})
				cmd.SetOutput(&output)

				Expect(cmd.Execute()).To(MatchError(fmt.Sprintf("failed to remove %s: some-error", path)))
			})
		})
	})

	Context("when the path has nested credentials", func() {
		Context("when `recursive` is true", func() {
			It("deletes all of the credentials", func() {
				rootPath := "/some/path"
				path1 := rootPath + "/some-name"
				path2 := rootPath + "/some-other-name"
				fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, &credhub.ErrCredentialNotFound{})
				fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{{Name: path1}, {Name: path2}}, nil)

				var output bytes.Buffer
				cmd := rm.NewCmdRm(dependencies)
				cmd.SetArgs([]string{"-r", rootPath})
				cmd.SetOutput(&output)

				Expect(cmd.Execute()).To(Succeed())

				Expect(fakeCredhubClient.GetCredentialByNameCallCount()).To(Equal(1))
				Expect(fakeCredhubClient.GetCredentialByNameArgsForCall(0)).To(Equal(rootPath))
				Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
				Expect(fakeCredhubClient.FindCredentialsByPathArgsForCall(0)).To(Equal(rootPath))
				Expect(fakeCredhubClient.DeleteCredentialByNameCallCount()).To(Equal(2))
				Expect(fakeCredhubClient.DeleteCredentialByNameArgsForCall(0)).To(Equal(path1))
				Expect(fakeCredhubClient.DeleteCredentialByNameArgsForCall(1)).To(Equal(path2))
			})

			Context("when deleting a credential fails", func() {
				It("returns an error", func() {
					path := "/path/to/some/cred"
					fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, &credhub.ErrCredentialNotFound{})
					fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{{Name: path}}, nil)
					fakeCredhubClient.DeleteCredentialByNameReturns(errors.New("some-error"))

					var output bytes.Buffer
					cmd := rm.NewCmdRm(dependencies)
					cmd.SetArgs([]string{"-r", "/"})
					cmd.SetOutput(&output)

					Expect(cmd.Execute()).To(MatchError(fmt.Sprintf("failed to remove %s: some-error", path)))
				})
			})

			Context("when finding credentials fails", func() {
				It("returns an error", func() {
					fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, &credhub.ErrCredentialNotFound{})
					fakeCredhubClient.FindCredentialsByPathReturns(nil, errors.New("some-error"))

					var output bytes.Buffer
					cmd := rm.NewCmdRm(dependencies)
					cmd.SetArgs([]string{"-r", "/"})
					cmd.SetOutput(&output)

					Expect(cmd.Execute()).To(MatchError("failed to find credentials: some-error"))
				})
			})
		})

		Context("when `recursive` is false", func() {
			It("returns an error", func() {
				rootPath := "/some/path"
				path1 := rootPath + "/some-name"
				path2 := rootPath + "/some-other-name"
				fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, &credhub.ErrCredentialNotFound{})
				fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{{Name: path1}, {Name: path2}}, nil)

				var output bytes.Buffer
				cmd := rm.NewCmdRm(dependencies)
				cmd.SetArgs([]string{rootPath})
				cmd.SetOutput(&output)

				Expect(cmd.Execute()).To(MatchError(errors.New("not removing recursively without '-r' flag")))
			})
		})
	})

	Context("when no arguments are provided", func() {
		It("returns an error and shows the usage", func() {
			cmd := rm.NewCmdRm(dependencies)
			cmd.SetArgs([]string{})
			cmd.SetOutput(ioutil.Discard)

			Expect(cmd.Execute()).To(MatchError(ContainSubstring("must provide a credential path")))
			Expect(cmd.SilenceUsage).To(BeFalse())
		})
	})

	Context("when the path matches no credentials", func() {
		It("returns an error", func() {
			fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, &credhub.ErrCredentialNotFound{})
			fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{}, nil)

			var output bytes.Buffer
			cmd := rm.NewCmdRm(dependencies)
			cmd.SetArgs([]string{"/some-path"})
			cmd.SetOutput(&output)

			Expect(cmd.Execute()).To(MatchError("'/some-path': no such credential or path"))
			Expect(cmd.SilenceUsage).To(BeTrue())
		})
	})

	Context("when getting the credential fails", func() {
		It("returns an error", func() {
			path := "/some/path/to/cred"
			fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, errors.New("some-error"))

			var output bytes.Buffer
			cmd := rm.NewCmdRm(dependencies)
			cmd.SetArgs([]string{path})
			cmd.SetOutput(&output)

			Expect(cmd.Execute()).To(MatchError("failed to get credential: some-error"))
		})
	})
})
