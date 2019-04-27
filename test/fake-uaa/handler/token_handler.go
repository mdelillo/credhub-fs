package handler

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *uaaHandler) tokenHandler(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	if grantType != "client_credentials" {
		panic("Only 'client_credentials' grant type is supported")
	}
	if clientID == "" {
		panic("'client_id' must not be empty")
	}
	if clientSecret == "" {
		panic("'client_secret' must not be empty")
	}

	if h.clients[clientID] != clientSecret {
		panic("'client_id' and/or 'client_secret' is incorrect")
	}

	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(h.jwtSigningKey))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse JWT signing key: %s", err.Error()))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"client_id":  clientID,
		"grant_type": "client_credentials",
		"iss":        fmt.Sprintf("https://%s%s", h.listenAddr, c.Request.URL.Path),
		"scope":      []string{"credhub.read", "credhub.write"},
	})
	token.Header["kid"] = "legacy-token-key"
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		panic(fmt.Sprintf("failed to sign JWT token: %s", err.Error()))
	}

	c.JSON(200, gin.H{"access_token": tokenString})
}
