package handler

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type uaaHandler struct {
	listenAddr    string
	jwtSigningKey *rsa.PrivateKey
	clients       map[string]string
}

func NewUAAHandler(listenAddr, jwtSigningKey string, clients []string) (http.Handler, error) {
	clientMap := make(map[string]string)
	for _, client := range clients {
		if strings.Count(client, ":") == 0 {
			return nil, errors.New("'clients' must contain colon-separated client IDs and secrets")
		}
		clientID := strings.SplitN(client, ":", 2)[0]
		clientSecret := strings.SplitN(client, ":", 2)[1]
		clientMap[clientID] = clientSecret
	}

	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(jwtSigningKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT signing key: %s", err.Error())
	}

	h := &uaaHandler{
		listenAddr:    listenAddr,
		jwtSigningKey: signingKey,
		clients:       clientMap,
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/oauth/token", h.tokenHandler)
	return router, nil
}
