package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *credhubHandler) putDataHandler(c *gin.Context) {
	var requestBody struct {
		Name  string `json:"name" binding:"required"`
		Type  string `json:"type" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		panic(fmt.Sprintf("Failed to parse request body: %s", err))
	}

	if requestBody.Type != "value" {
		panic("Only 'value' types are supported")
	}
	if requestBody.Name == "" {
		panic("'name' must not be empty")
	}
	if requestBody.Value == "" {
		panic("'value' must not be empty")
	}

	createdCred := Credential{
		ID:               uuid.New(),
		VersionCreatedAt: time.Now().UTC(),
		Name:             requestBody.Name,
		Value:            requestBody.Value,
		Type:             requestBody.Type,
	}

	h.credentials[requestBody.Name] = createdCred

	c.JSON(200, createdCred)
}
