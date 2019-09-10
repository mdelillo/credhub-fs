package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *credhubHandler) deleteDataHandler(c *gin.Context) {
	name := c.Query("name")

	if name == "" {
		c.JSON(400, gin.H{
			"error": ErrMissingNameParameter,
		})
	} else {
		found := h.credentialStore.Delete(name)
		if !found {
			c.JSON(404, gin.H{
				"error": ErrCredentialDoesNotExist,
			})
			return
		}

		c.Status(204)
	}
}
