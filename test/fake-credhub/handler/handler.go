package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type credhubHandler struct {
	authServerURL      string
	jwtVerificationKey string
	credentials        map[string]Credential
}

type Credential struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Value            string    `json:"value"`
	VersionCreatedAt time.Time `json:"version_created_at"`
}

func NewCredhubHandler(authServerURL, jwtVerificationKey string) http.Handler {
	h := &credhubHandler{
		authServerURL:      authServerURL,
		jwtVerificationKey: jwtVerificationKey,
		credentials:        make(map[string]Credential),
	}

	router := gin.Default()
	router.GET("/info", h.infoHandler)

	authenticationRequired := router.Group("/", h.authenticationRequired)
	{
		authenticationRequired.GET("/api/v1/data", h.getDataHandler)
		authenticationRequired.PUT("/api/v1/data", h.putDataHandler)
	}

	return router
}
