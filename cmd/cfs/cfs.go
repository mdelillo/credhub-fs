package main

import (
	"fmt"
	"os"

	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd"
)

func main() {
	command := cmd.NewCfsCommand()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
