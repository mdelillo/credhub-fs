package cmd

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/mdelillo/credhub-fs/pkg/cfs/cmd/cat"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCfsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cfs",
		Short: "cfs interacts with CredHub using Unix filesystem commands",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	viper.AutomaticEnv()

	var (
		credhubAddr  string
		clientID     string
		clientSecret string
	)
	cmd.PersistentFlags().StringVarP(&credhubAddr, "credhub-addr", "", viper.GetString("CREDHUB_ADDR"), "Address of CredHub server [$CREDHUB_ADDR]")
	cmd.PersistentFlags().StringVarP(&clientID, "client-id", "", viper.GetString("CLIENT_ID"), "UAA client ID [$CLIENT_ID]")
	cmd.PersistentFlags().StringVarP(&clientSecret, "client-secret", "", viper.GetString("CLIENT_SECRET"), "UAA client secret [$CLIENT_SECRET]")
	cmd.MarkFlagRequired("credhub-addr")
	cmd.MarkFlagRequired("client-id")
	cmd.MarkFlagRequired("client-secret")

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	cmd.AddCommand(cat.NewCmdCat(httpClient))

	return cmd
}
