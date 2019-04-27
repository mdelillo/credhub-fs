package cat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

type cmdCatRunner struct {
	httpClient *http.Client
}

func NewCmdCat(httpClient *http.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "cat",
		Short: "Get the value of a credential",
		Run: func(cmd *cobra.Command, args []string) {
			c := &cmdCatRunner{httpClient: httpClient}
			if err := c.Run(cmd, args); err != nil {
				log.Fatal(err)
			}
		},
	}
}

func (c *cmdCatRunner) Run(cmd *cobra.Command, args []string) error {
	credhubAddr, err := cmd.Flags().GetString("credhub-addr")
	if err != nil {
		log.Fatal(err)
	}
	clientID, err := cmd.Flags().GetString("client-id")
	if err != nil {
		log.Fatal(err)
	}
	clientSecret, err := cmd.Flags().GetString("client-secret")
	if err != nil {
		log.Fatal(err)
	}

	authToken, err := c.getUAAToken(credhubAddr, clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	name := args[0]
	credentialValue, err := c.getCredentialValue(name, credhubAddr, authToken)
	if err != nil {
		log.Fatal(err)
	}

	if credentialValue != "" {
		fmt.Println(credentialValue)
	}
	return nil
}

func (c *cmdCatRunner) getCredentialValue(name, credhubAddr, authToken string) (string, error) {
	url := fmt.Sprintf("https://%s/api/v1/data?name=%s", credhubAddr, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var credentials struct {
		Data []struct {
			Value string `json:"value"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		log.Fatal(err)
	}

	if len(credentials.Data) == 0 {
		return "", nil
	}

	if len(credentials.Data) > 1 {
		log.Fatalf("Expected one credential but got %d", len(credentials.Data))
	}

	return credentials.Data[0].Value, nil
}

func (c *cmdCatRunner) getUAAURL(credhubAddr string) (string, error) {
	infoURL := fmt.Sprintf("https://%s/info", credhubAddr)
	resp, err := c.httpClient.Get(infoURL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var infoResponse struct {
		AuthServer struct {
			URL string `json:"url"`
		} `json:"auth-server"`
	}
	if err := json.Unmarshal(body, &infoResponse); err != nil {
		log.Fatal(err)
	}

	return infoResponse.AuthServer.URL, nil
}

func (c *cmdCatRunner) getUAAToken(credhubAddr, clientID, clientSecret string) (string, error) {
	uaaURL, err := c.getUAAURL(credhubAddr)
	if err != nil {
		log.Fatal(err)
	}

	tokenURL := fmt.Sprintf("%s/oauth/token", uaaURL)
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"grant_type":    {"client_credentials"},
	}
	resp, err := c.httpClient.PostForm(tokenURL, values)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		log.Fatal(err)
	}

	return tokenResponse.AccessToken, nil
}
