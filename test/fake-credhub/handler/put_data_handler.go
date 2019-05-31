package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/credentials"
)

func (h *credhubHandler) putDataHandler(c *gin.Context) {
	var requestBody struct {
		Name  string `json:"name" binding:"required"`
		Type  string `json:"type" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{
			"error": ErrInvalidPathOrBody,
		})
		return
	}

	if requestBody.Type != "value" {
		c.JSON(400, gin.H{
			"error": ErrInvalidType,
		})
		return
	}

	createdCred := credentials.Credential{
		ID:               uuid.New(),
		VersionCreatedAt: time.Now().UTC(),
		Name:             requestBody.Name,
		Value:            requestBody.Value,
		Type:             requestBody.Type,
	}

	h.credentialStore.Set(createdCred)

	c.JSON(200, createdCred)
}
