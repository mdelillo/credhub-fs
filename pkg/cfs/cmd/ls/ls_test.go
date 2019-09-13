package ls_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/ls"
	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/ls/lsfakes"
	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
)

var _ = Describe("Ls", func() {
	var fakeCredhubClient *lsfakes.FakeCredhubClient
	var dependencies cmdutil.Dependencies

	BeforeEach(func() {
		fakeCredhubClient = &lsfakes.FakeCredhubClient{}
		dependencies = cmdutil.NewDependencies()
		dependencies.SetCredhubClient(fakeCredhubClient)
	})

	It("lists top-level credentials and directories, sorted and uniqued", func() {
		credentials := []credhub.Credential{
			{Name: "/some-dir/cred1"},
			{Name: "/some-cred"},
			{Name: "/some-dir/cred2"},
			{Name: "/some-other-dir/some-nested-dir/cred"},
		}
		fakeCredhubClient.FindCredentialsByPathReturns(credentials, nil)

		var output bytes.Buffer
		cmd := ls.NewCmdLs(dependencies)
		cmd.SetOutput(&output)
		cmd.SetArgs([]string{})

		Expect(cmd.Execute()).To(Succeed())

		Expect(output.String()).To(Equal("/some-cred  /some-dir/  /some-other-dir/\n"))
		Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
		Expect(fakeCredhubClient.FindCredentialsByPathArgsForCall(0)).To(Equal("/"))
	})

	Context("when a path is specified", func() {
		It("lists one level of credentials and directories at that path, sorted and uniqued", func() {
			path := "/some-dir"
			credentials := []credhub.Credential{
				{Name: "/some-dir/some-nested-dir/cred"},
				{Name: "/some-dir/cred1"},
				{Name: "/some-dir/cred2"},
			}
			fakeCredhubClient.FindCredentialsByPathReturns(credentials, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{path})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(Equal("/some-dir/cred1  /some-dir/cred2  /some-dir/some-nested-dir/\n"))
			Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.FindCredentialsByPathArgsForCall(0)).To(Equal(path))
		})
	})

	Context("when a path is specified with a trailing slash", func() {
		It("lists one level of credentials and directories at that path, sorted and uniqued", func() {
			path := "/some-dir/"
			credentials := []credhub.Credential{
				{Name: "/some-dir/some-nested-dir/cred"},
				{Name: "/some-dir/cred1"},
				{Name: "/some-dir/cred2"},
			}
			fakeCredhubClient.FindCredentialsByPathReturns(credentials, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{path})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(Equal("/some-dir/cred1  /some-dir/cred2  /some-dir/some-nested-dir/\n"))
			Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.FindCredentialsByPathArgsForCall(0)).To(Equal(path))
		})
	})

	Context("when a path is specified without a leading slash", func() {
		It("lists one level of credentials and directories at that path, sorted and uniqued", func() {
			path := "some-dir/"
			credentials := []credhub.Credential{
				{Name: "/some-dir/some-nested-dir/cred"},
				{Name: "/some-dir/cred1"},
				{Name: "/some-dir/cred2"},
			}
			fakeCredhubClient.FindCredentialsByPathReturns(credentials, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{path})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(Equal("/some-dir/cred1  /some-dir/cred2  /some-dir/some-nested-dir/\n"))
			Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.FindCredentialsByPathArgsForCall(0)).To(Equal("/" + path))
		})
	})

	Context("when the path matches a credential exactly", func() {
		It("should list the credential", func() {
			path := "/some-cred"
			credential := credhub.Credential{Name: path}
			fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{}, nil)
			fakeCredhubClient.GetCredentialByNameReturns(credential, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{path})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(Equal("/some-cred\n"))

			Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.FindCredentialsByPathArgsForCall(0)).To(Equal(path))
			Expect(fakeCredhubClient.GetCredentialByNameCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.GetCredentialByNameArgsForCall(0)).To(Equal(path))
		})
	})

	Context("when the path matches a nested credential exactly", func() {
		It("should list the credential", func() {
			path := "/some-dir/some-cred"
			credential := credhub.Credential{Name: path}
			fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{}, nil)
			fakeCredhubClient.GetCredentialByNameReturns(credential, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{path})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(Equal("/some-dir/some-cred\n"))

			Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.FindCredentialsByPathArgsForCall(0)).To(Equal(path))
			Expect(fakeCredhubClient.GetCredentialByNameCallCount()).To(Equal(1))
			Expect(fakeCredhubClient.GetCredentialByNameArgsForCall(0)).To(Equal(path))
		})
	})

	Context("when no credentials are found", func() {
		Context("when the path is '/'", func() {
			It("does not print anything", func() {
				fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{}, nil)

				var output bytes.Buffer
				cmd := ls.NewCmdLs(dependencies)
				cmd.SetOutput(&output)
				cmd.SetArgs([]string{})

				Expect(cmd.Execute()).To(Succeed())

				Expect(output.String()).To(Equal("\n"))
			})
		})

		Context("when the path is not '/'", func() {
			It("returns an error", func() {
				fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{}, nil)
				fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, &credhub.ErrCredentialNotFound{})

				var output bytes.Buffer
				cmd := ls.NewCmdLs(dependencies)
				cmd.SetOutput(&output)
				cmd.SetArgs([]string{"/some-path"})

				Expect(cmd.Execute()).To(MatchError("'/some-path': no such credential or path"))
				Expect(cmd.SilenceUsage).To(BeTrue())

				Expect(fakeCredhubClient.FindCredentialsByPathCallCount()).To(Equal(1))
				Expect(fakeCredhubClient.GetCredentialByNameCallCount()).To(Equal(1))
			})
		})
	})

	Context("when the '1' option is specified", func() {
		It("lists one credential per line", func() {
			credentials := []credhub.Credential{
				{Name: "/some-cred-1"},
				{Name: "/some-cred-2"},
			}
			fakeCredhubClient.FindCredentialsByPathReturns(credentials, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{"-1"})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(Equal("/some-cred-1\n/some-cred-2\n"))
		})
	})

	Context("when the 'long' option is specified", func() {
		It("shows one credential name and date per line", func() {
			credentials := []credhub.Credential{
				{Name: "/some-cred-1", VersionCreatedAt: time.Date(1985, time.October, 26, 0, 0, 0, 0, time.UTC)},
				{Name: "/some-cred-2", VersionCreatedAt: time.Date(time.Now().Year(), time.March, 1, 17, 4, 30, 0, time.UTC)},
			}
			fakeCredhubClient.FindCredentialsByPathReturns(credentials, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{"-l"})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(ContainSubstring("Oct 26  1985 /some-cred-1"))
			Expect(output.String()).To(ContainSubstring("Mar  1 17:04 /some-cred-2"))
		})
	})

	Context("when both the '1' and 'long' options are specified", func() {
		It("prefers 'long'", func() {
			credentials := []credhub.Credential{
				{Name: "/some-cred-1", VersionCreatedAt: time.Date(1985, time.October, 26, 0, 0, 0, 0, time.UTC)},
				{Name: "/some-cred-2", VersionCreatedAt: time.Date(time.Now().Year(), time.March, 1, 7, 4, 30, 0, time.UTC)},
				{Name: "/some-cred-3", VersionCreatedAt: time.Date(time.Now().Year(), time.March, 11, 17, 14, 30, 0, time.UTC)},
			}
			fakeCredhubClient.FindCredentialsByPathReturns(credentials, nil)

			var output bytes.Buffer
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(&output)
			cmd.SetArgs([]string{"-l", "-1"})

			Expect(cmd.Execute()).To(Succeed())

			Expect(output.String()).To(ContainSubstring("Oct 26  1985 /some-cred-1"))
			Expect(output.String()).To(ContainSubstring("Mar  1 07:04 /some-cred-2"))
			Expect(output.String()).To(ContainSubstring("Mar 11 17:14 /some-cred-3"))
		})
	})

	Context("when the '1' option is specified incorrectly", func() {
		It("returns an error and prints the usage", func() {
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(ioutil.Discard)
			cmd.SetArgs([]string{"-1=blah"})

			Expect(cmd.Execute()).To(MatchError(ContainSubstring("invalid argument")))
			Expect(cmd.SilenceUsage).To(BeFalse())
		})
	})

	Context("when the 'l' option is specified incorrectly", func() {
		It("returns an error and prints the usage", func() {
			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(ioutil.Discard)
			cmd.SetArgs([]string{"-l=blah"})

			Expect(cmd.Execute()).To(MatchError(ContainSubstring("invalid argument")))
			Expect(cmd.SilenceUsage).To(BeFalse())
		})
	})

	Context("when finding credentials fails", func() {
		It("returns an error", func() {
			fakeCredhubClient.FindCredentialsByPathReturns(nil, errors.New("some-error"))

			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(ioutil.Discard)
			cmd.SetArgs([]string{})

			Expect(cmd.Execute()).To(MatchError("failed to list credentials: some-error"))
			Expect(cmd.SilenceUsage).To(BeTrue())
		})
	})

	Context("when getting a credential by name fails", func() {
		It("returns an error", func() {
			fakeCredhubClient.FindCredentialsByPathReturns([]credhub.Credential{}, nil)
			fakeCredhubClient.GetCredentialByNameReturns(credhub.Credential{}, errors.New("some-error"))

			cmd := ls.NewCmdLs(dependencies)
			cmd.SetOutput(ioutil.Discard)
			cmd.SetArgs([]string{"/some-path"})

			Expect(cmd.Execute()).To(MatchError("failed to get credential: some-error"))
			Expect(cmd.SilenceUsage).To(BeTrue())
		})
	})
})
