package handler

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h *uaaHandler) tokenHandler(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	if grantType != "client_credentials" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Grant type must be 'client_credentials'"})
		return
	}
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'client_id' must not be empty"})
		return
	}
	if clientSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'client_secret' must not be empty"})
		return
	}

	if h.clients[clientID] != clientSecret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect 'client_id' and/or 'client_secret'"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"client_id":  clientID,
		"grant_type": "client_credentials",
		"iss":        fmt.Sprintf("https://%s%s", h.listenAddr, c.Request.URL.Path),
		"scope":      []string{"credhub.read", "credhub.write"},
	})
	token.Header["kid"] = "legacy-token-key"
	tokenString, err := token.SignedString(h.jwtSigningKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to sign JWT token: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokenString})
}
