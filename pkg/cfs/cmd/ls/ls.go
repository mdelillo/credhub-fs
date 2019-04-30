package ls

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
	"github.com/spf13/cobra"
)

type cmdLsRunner struct {
	credhubClient credhubClient
	formatLong    bool
	formatOne     bool
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . credhubClient
type credhubClient credhub.Client

func NewCmdLs(dependencies cmdutil.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			formatLong, _ := cmd.Flags().GetBool("l")
			formatOne, _ := cmd.Flags().GetBool("1")

			cmd.SilenceUsage = true

			c := &cmdLsRunner{
				credhubClient: dependencies.GetCredhubClient(),
				formatLong:    formatLong,
				formatOne:     formatOne,
			}
			return c.Run(cmd, args)
		},
	}

	cmd.Flags().BoolP("l", "l", false, "list in long format")
	cmd.Flags().BoolP("1", "1", false, "list one per line")

	return cmd
}

func (c *cmdLsRunner) Run(cmd *cobra.Command, args []string) error {
	path := "/"
	if len(args) > 0 {
		path = args[0]
	}

	credentials, err := c.credhubClient.FindCredentialsByPath(path)
	if err != nil {
		return fmt.Errorf("failed to list credentials: %s", err.Error())
	}

	if len(credentials) == 0 && path != "/" {
		credential, err := c.credhubClient.GetCredentialByName(path)
		if err != nil {
			switch err.(type) {
			case *credhub.ErrCredentialNotFound:
				return fmt.Errorf("'%s': no such credential or path", path)
			default:
				return fmt.Errorf("failed to get credential: %s", err.Error())
			}
		}
		credentials = []credhub.Credential{credential}
	}

	output := c.formatCredentialPathsAndDates(credentials, path)
	output = c.sort(output)
	output = c.uniq(output)

	separator := "  "
	if c.formatLong || c.formatOne {
		separator = "\n"
	}
	fmt.Fprintln(cmd.OutOrStdout(), strings.Join(output, separator))
	return nil
}

func (c *cmdLsRunner) formatCredentialPathsAndDates(credentials []credhub.Credential, path string) []string {
	var output []string
	for _, credential := range credentials {
		credentialOutput := credential.Name
		if strings.Count(credentialOutput, "/") > 1 {
			name := strings.TrimPrefix(credentialOutput, strings.TrimSuffix(path, "/"))
			credentialOutput = filepath.Join(path, strings.Split(name, "/")[1])
			if strings.Count(name, "/") > 1 {
				credentialOutput = credentialOutput + "/"
			}
		}

		if c.formatLong {
			credentialOutput = c.prependDate(credentialOutput, credential.VersionCreatedAt)
		}

		output = append(output, credentialOutput)
	}
	return output
}

func (c *cmdLsRunner) sort(a []string) []string {
	sort.Strings(a)
	return a
}

func (c *cmdLsRunner) uniq(a []string) []string {
	seen := make(map[string]struct{}, len(a))
	count := 0
	for _, v := range a {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		a[count] = v
		count++
	}
	return a[:count]
}

func (c *cmdLsRunner) prependDate(credentialOutput string, date time.Time) string {
	var dateString string
	if date.Year() == time.Now().Year() {
		dateString = date.Format("Jan _2 15:04")
	} else {
		dateString = date.Format("Jan _2  2006")
	}
	return dateString + " " + credentialOutput
}
