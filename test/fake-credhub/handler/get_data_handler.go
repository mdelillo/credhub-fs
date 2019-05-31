package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mdelillo/credhub-fs/test/fake-credhub/credentials"
)

func (h *credhubHandler) getDataHandler(c *gin.Context) {
	name := c.Query("name")
	path := c.Query("path")

	if name == "" && path == "" {
		c.JSON(400, gin.H{
			"error": ErrMissingNameParameter,
		})
	} else if path != "" {
		h.getDataByPathHandler(path, c)
	} else {
		h.getDataByNameHandler(name, c)
	}
}

func (h *credhubHandler) getDataByNameHandler(name string, c *gin.Context) {
	cred, found := h.credentialStore.GetByName(name)
	if !found {
		c.JSON(404, gin.H{
			"error": ErrCredentialDoesNotExist,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": []credentials.Credential{cred},
	})
}

func (h *credhubHandler) getDataByPathHandler(path string, c *gin.Context) {
	creds := h.credentialStore.GetByPath(path)
	var credsView []credentials.CredentialNameAndDate
	for _, cred := range creds {
		credsView = append(credsView, credentials.CredentialNameAndDate{
			Name: cred.Name, VersionCreatedAt: cred.VersionCreatedAt,
		})
	}

	c.JSON(200, gin.H{
		"credentials": credsView,
	})
}
