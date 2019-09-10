package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/cat"
	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/ls"
	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/rm"
	cmdutil "github.com/mdelillo/credhub-fs/pkg/cfs/cmd/util"
	"github.com/mdelillo/credhub-fs/pkg/credhub"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCfsCommand() *cobra.Command {
	dependencies := cmdutil.NewDependencies()
	cmd := &cobra.Command{
		Use:   "cfs",
		Short: "cfs interacts with CredHub using Unix filesystem commands",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			requiredFlags := []string{"credhub-addr", "client-id", "client-secret"}
			for _, flag := range requiredFlags {
				if viper.GetString(flag) == "" {
					fmt.Printf("Must provide `%s`\n", flag)
					cmd.Usage()
					os.Exit(1)
				}
			}

			httpClient := &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
					Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
					TLSHandshakeTimeout: 5 * time.Second,
				},
			}

			dependencies.SetCredhubClient(
				credhub.NewClient(
					viper.GetString("credhub-addr"),
					viper.GetString("client-id"),
					viper.GetString("client-secret"),
					httpClient,
				),
			)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.PersistentFlags().String("credhub-addr", "", "address of CredHub server [$CREDHUB_ADDR]")
	cmd.PersistentFlags().String("client-id", "", "UAA client ID [$CLIENT_ID]")
	cmd.PersistentFlags().String("client-secret", "", "UAA client secret [$CLIENT_SECRET]")
	viper.BindEnv("credhub-addr", "CREDHUB_ADDR")
	viper.BindEnv("client-id", "CLIENT_ID")
	viper.BindEnv("client-secret", "CLIENT_SECRET")
	viper.BindPFlags(cmd.PersistentFlags())

	cmd.AddCommand(cat.NewCmdCat(dependencies))
	cmd.AddCommand(ls.NewCmdLs(dependencies))
	cmd.AddCommand(rm.NewCmdRm(dependencies))

	return cmd
}
