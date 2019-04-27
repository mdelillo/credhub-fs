package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *credhubHandler) infoHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"auth-server": gin.H{
			"url": h.authServerURL,
		},
		"app": gin.H{
			"name": "Fake CredHub",
		},
	})
}
