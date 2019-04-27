package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *credhubHandler) getDataHandler(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		panic("Must provide 'name'")
	}

	cred, credExists := h.credentials[name]
	if !credExists {
		panic("Could not find cred " + name)
	}

	c.JSON(200, gin.H{
		"data": []Credential{cred},
	})
}
