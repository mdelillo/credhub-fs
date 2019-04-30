package cat

import (
	"errors"
	"fmt"

	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
	"github.com/spf13/cobra"
)

type cmdCatRunner struct {
	credhubClient credhubClient
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . credhubClient
type credhubClient credhub.Client

func NewCmdCat(dependencies cmdutil.Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "cat /path/to/credential",
		Short: "Get the value of a credential",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must provide a credential path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			c := &cmdCatRunner{credhubClient: dependencies.GetCredhubClient()}
			return c.Run(cmd, args)
		},
	}
}

func (c *cmdCatRunner) Run(cmd *cobra.Command, args []string) error {
	name := args[0]
	cred, err := c.credhubClient.GetCredentialByName(name)
	if err != nil {
		switch err.(type) {
		case *credhub.ErrCredentialNotFound:
			return fmt.Errorf("'%s': no such credential or path", name)
		default:
			return fmt.Errorf("failed to get credential: %s", err.Error())
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), cred.Value)
	return nil
}
