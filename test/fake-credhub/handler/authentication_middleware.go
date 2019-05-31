package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *credhubHandler) authenticationRequired(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		c.JSON(401, gin.H{
			"error":             ErrInvalidToken,
			"error_description": ErrDescriptionNoAuthentication,
		})
		c.Abort()
		return
	}

	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	claims := map[string]string{"iss": h.authServerURL}
	if err := h.tokenValidator.ValidateTokenWithClaims(tokenString, claims); err != nil {
		c.JSON(401, gin.H{
			"error":             ErrInvalidToken,
			"error_description": ErrDescriptionMalformedToken,
		})
		c.Abort()
		return
	}
}
