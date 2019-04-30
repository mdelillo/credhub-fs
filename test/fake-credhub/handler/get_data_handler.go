package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *credhubHandler) getDataHandler(c *gin.Context) {
	name := c.Query("name")
	path := c.Query("path")

	if name == "" && path == "" {
		panic("Must provide 'name' or 'path'")
	} else if name != "" && path != "" {
		panic("Cannot provide both 'name' and 'path'")
	} else if name != "" {
		h.getDataByNameHandler(name, c)
	} else {
		h.getDataByPathHandler(path, c)
	}
}

func (h *credhubHandler) getDataByNameHandler(name string, c *gin.Context) {
	cred, credExists := h.credentials[name]
	if !credExists {
		panic("Could not find cred " + name)
	}

	c.JSON(200, gin.H{
		"data": []Credential{cred},
	})
}

func (h *credhubHandler) getDataByPathHandler(path string, c *gin.Context) {
	type matchingCred struct {
		Name             string    `json:"name"`
		VersionCreatedAt time.Time `json:"version_created_at"`
	}
	var matchingCreds []matchingCred

	for name, cred := range h.credentials {
		if strings.HasPrefix(name, path) && name != path {
			matchingCreds = append(matchingCreds, matchingCred{
				Name:             cred.Name,
				VersionCreatedAt: cred.VersionCreatedAt,
			})
		}
	}

	c.JSON(200, gin.H{
		"credentials": matchingCreds,
	})
}
