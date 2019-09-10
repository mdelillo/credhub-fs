package rm

import (
	"errors"
	"fmt"

	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
	"github.com/spf13/cobra"
)

type cmdRmRunner struct {
	credhubClient credhubClient
	recursive     bool
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . credhubClient
type credhubClient credhub.Client

func NewCmdRm(dependencies cmdutil.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm /path/to/credential",
		Short: "Removes a credential",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must provide a credential path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			recursive, _ := cmd.Flags().GetBool("recursive")

			cmd.SilenceUsage = true

			c := &cmdRmRunner{
				credhubClient: dependencies.GetCredhubClient(),
				recursive:     recursive,
			}
			return c.Run(cmd, args)
		},
	}

	cmd.Flags().BoolP("recursive", "r", false, "recursively delete credentials")

	return cmd
}

func (c *cmdRmRunner) Run(cmd *cobra.Command, args []string) error {
	path := args[0]
	credential, err := c.credhubClient.GetCredentialByName(path)
	if err != nil {
		if _, isNotFoundError := err.(*credhub.ErrCredentialNotFound); !isNotFoundError {
			return fmt.Errorf("failed to get credential: %s", err.Error())
		}
	}

	if credential.Name != "" {
		if err := c.credhubClient.DeleteCredentialByName(path); err != nil {
			return fmt.Errorf("failed to remove %s: %s", path, err.Error())
		}
	} else {
		credentials, err := c.credhubClient.FindCredentialsByPath(path)
		if err != nil {
			return fmt.Errorf("failed to find credentials: %s", err.Error())
		}

		if len(credentials) > 0 {
			if c.recursive {
				for _, credential := range credentials {
					if err := c.credhubClient.DeleteCredentialByName(credential.Name); err != nil {
						return fmt.Errorf("failed to remove %s: %s", credential.Name, err.Error())
					}
				}
			} else {
				return errors.New("not removing recursively without '-r' flag")
			}
		} else {
			return fmt.Errorf("'%s': no such credential or path", path)
		}
	}

	return nil
}
