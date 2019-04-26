package cmd

import "github.com/spf13/cobra"

func NewCfsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cfs",
		Short: "cfs interacts with CredHub using Unix filesystem commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}
