package main

import (
	"os"

	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd"
)

func main() {
	command := cmd.NewCfsCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
