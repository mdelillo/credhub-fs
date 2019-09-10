package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/credentials"
)

type credhubHandler struct {
	authServerURL   string
	credentialStore credentialStore
	tokenValidator  tokenValidator
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . credentialStore
type credentialStore interface {
	GetByName(name string) (cred credentials.Credential, found bool)
	GetByPath(path string) []credentials.Credential
	Set(credential credentials.Credential)
	Delete(name string) bool
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . tokenValidator
type tokenValidator interface {
	ValidateTokenWithClaims(token string, claims map[string]string) error
}

func NewCredhubHandler(authServerURL string, credentialStore credentialStore, tokenValidator tokenValidator) (http.Handler, error) {
	h := &credhubHandler{
		authServerURL:   authServerURL,
		credentialStore: credentialStore,
		tokenValidator:  tokenValidator,
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/info", h.infoHandler)

	authenticationRequired := router.Group("/", h.authenticationRequired)
	{
		authenticationRequired.GET("/api/v1/data", h.getDataHandler)
		authenticationRequired.PUT("/api/v1/data", h.putDataHandler)
		authenticationRequired.DELETE("/api/v1/data", h.deleteDataHandler)
	}

	return router, nil
}
