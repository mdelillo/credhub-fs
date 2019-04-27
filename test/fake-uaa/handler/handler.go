package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type uaaHandler struct {
	listenAddr    string
	jwtSigningKey string
	clients       map[string]string
}

func NewUAAHandler(listenAddr, jwtSigningKey string, clients []string) http.Handler {
	clientMap := make(map[string]string)
	for _, client := range clients {
		clientID := strings.Split(client, ":")[0]
		clientSecret := strings.Split(client, ":")[1]
		clientMap[clientID] = clientSecret
	}

	h := &uaaHandler{
		listenAddr:    listenAddr,
		jwtSigningKey: jwtSigningKey,
		clients:       clientMap,
	}

	router := gin.Default()
	router.POST("/oauth/token", h.tokenHandler)
	return router
}
