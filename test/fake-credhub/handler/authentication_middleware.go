package handler

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	}

	tokenString := strings.TrimPrefix(authorization, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, signedUsingRSA := token.Method.(*jwt.SigningMethodRSA); !signedUsingRSA {
			panic("Token not signed using RSA")
		}

		jwtVerificationKey, err := h.parseJWTVerificationKey()
		if err != nil {
			panic(err)
		}

		return jwtVerificationKey, nil
	})
	if err != nil {
		panic(err)
	}

	if !token.Valid {
		panic("Token is invalid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if !strings.HasPrefix(claims["iss"].(string), h.authServerURL) {
			c.JSON(401, gin.H{
				"error":             ErrInvalidToken,
				"error_description": ErrDescriptionMalformedToken,
			})
			c.Abort()
		}
	} else {
		panic("Failed to parse claims")
	}
}

func (h *credhubHandler) parseJWTVerificationKey() (*rsa.PublicKey, error) {
	key := strings.TrimSpace(h.jwtVerificationKey)
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey.(*rsa.PublicKey), nil
}
