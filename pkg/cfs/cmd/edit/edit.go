package edit

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
	"github.com/spf13/cobra"
)

type cmdEditRunner struct {
	credhubClient credhubClient
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . credhubClient
type credhubClient credhub.Client

// TODO: add `--type` flag that defaults to value but can be JSON
func NewCmdEdit(dependencies cmdutil.Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:     "edit /path/to/credential",
		Short:   "Edit the value of a credential using $EDITOR",
		Aliases: []string{"vim"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("must provide a credential path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			c := &cmdEditRunner{credhubClient: dependencies.GetCredhubClient()}
			return c.Run(cmd, args)
		},
	}
}

func (c *cmdEditRunner) Run(cmd *cobra.Command, args []string) error {
	name := args[0]

	cred, err := c.credhubClient.GetCredentialByName(name)
	if err != nil {
		if _, isNotFoundError := err.(*credhub.ErrCredentialNotFound); !isNotFoundError {
			return fmt.Errorf("failed to get credential: %s", err.Error())
		}
	}

	var credType string
	if cred.Type == "" {
		credType = "value"
	} else if cred.Type == "value" || cred.Type == "json" {
		credType = cred.Type
	} else {
		panic("can only edit `value` and `json` types")
	}

	newCredValue, err := writeAndEditFile([]byte(cred.Value))
	if err != nil {
		log.Fatal(err)
	}

	if err := c.credhubClient.Set(name, credType, strings.TrimSpace(string(newCredValue))); err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "%s updated\n", name)
	return nil
}

func writeAndEditFile(value []byte) ([]byte, error) {
	file, err := ioutil.TempFile("", "cfs")
	if err != nil {
		log.Fatal(err)
	}
	fileName := file.Name()

	if _, err = file.Write(value); err != nil {
		log.Fatal(err)
	}
	if err = file.Close(); err != nil {
		log.Fatal(err)
	}

	editor, exists := os.LookupEnv("EDITOR")
	if !exists {
		vimPath, err := exec.LookPath("vim")
		if err != nil {
			log.Fatal(err)
		}
		editor = vimPath
	}

	edit := exec.Command(editor, fileName)
	edit.Stdin = os.Stdin
	edit.Stdout = os.Stdout
	edit.Stderr = os.Stderr

	if err := edit.Run(); err != nil {
		log.Fatal(err)
	}

	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return contents, err
}
