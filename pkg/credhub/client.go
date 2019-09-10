package credhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type client struct {
	credhubAddr  string
	clientID     string
	clientSecret string
	uaaURL       string
	httpClient   *http.Client
}

type Client interface {
	DeleteCredentialByName(name string) error
	GetCredentialByName(name string) (Credential, error)
	FindCredentialsByPath(path string) ([]Credential, error)
}

func NewClient(credhubAddr, clientID, clientSecret string, httpClient *http.Client) Client {
	return &client{
		credhubAddr:  credhubAddr,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   httpClient,
	}
}

func (c *client) DeleteCredentialByName(name string) error {
	url := fmt.Sprintf("https://%s/api/v1/data?name=%s", c.credhubAddr, name)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err.Error())
	}

	authToken, err := c.getToken()
	if err != nil {
		return fmt.Errorf("failed to get token: %s", err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %s", err.Error())
	}

	if resp.StatusCode == http.StatusNotFound {
		return &ErrCredentialNotFound{name}
	} else if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("got %s", resp.Status)
	}

	return nil
}

func (c *client) GetCredentialByName(name string) (Credential, error) {
	url := fmt.Sprintf("https://%s/api/v1/data?name=%s", c.credhubAddr, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to create request: %s", err.Error())
	}

	authToken, err := c.getToken()
	if err != nil {
		return Credential{}, fmt.Errorf("failed to get token: %s", err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to make request: %s", err.Error())
	}

	if resp.StatusCode == http.StatusNotFound {
		return Credential{}, &ErrCredentialNotFound{name}
	} else if resp.StatusCode != http.StatusOK {
		return Credential{}, fmt.Errorf("got %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Credential{}, fmt.Errorf("failed to read body: %s", err.Error())
	}

	var credentials struct {
		Data []Credential `json:"data"`
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		return Credential{}, fmt.Errorf("failed to parse response body: %s\n%s", err.Error(), string(body))
	}

	if len(credentials.Data) != 1 {
		return Credential{}, fmt.Errorf("expected 1 credential but got %d", len(credentials.Data))
	}

	return credentials.Data[0], nil
}

func (c *client) FindCredentialsByPath(path string) ([]Credential, error) {
	url := fmt.Sprintf("https://%s/api/v1/data?path=%s", c.credhubAddr, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err.Error())
	}

	authToken, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %s", err.Error())
	}

	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err.Error())
	}

	var credentials struct {
		Credentials []Credential `json:"credentials"`
	}

	if err := json.Unmarshal(body, &credentials); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %s\n%s", err.Error(), string(body))
	}

	return credentials.Credentials, nil
}

func (c *client) getToken() (string, error) {
	if c.uaaURL == "" {
		uaaURL, err := c.getUAAURL()
		if err != nil {
			return "", fmt.Errorf("failed to get UAA URL: %s", err.Error())
		}
		c.uaaURL = uaaURL
	}

	tokenURL := fmt.Sprintf("%s/oauth/token", c.uaaURL)
	values := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"grant_type":    {"client_credentials"},
	}
	resp, err := c.httpClient.PostForm(tokenURL, values)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %s", err.Error())
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to parse response body: %s\n%s", err.Error(), string(body))
	}

	return tokenResponse.AccessToken, nil
}

func (c *client) getUAAURL() (string, error) {
	infoURL := fmt.Sprintf("https://%s/info", c.credhubAddr)
	resp, err := c.httpClient.Get(infoURL)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %s", err.Error())
	}

	var infoResponse struct {
		AuthServer struct {
			URL string `json:"url"`
		} `json:"auth-server"`
	}
	if err := json.Unmarshal(body, &infoResponse); err != nil {
		return "", fmt.Errorf("failed to parse response body: %s\n%s", err.Error(), string(body))
	}

	return infoResponse.AuthServer.URL, nil
}
