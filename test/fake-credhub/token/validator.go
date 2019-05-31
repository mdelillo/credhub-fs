package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type validator struct {
	jwtVerificationKey *rsa.PublicKey
}

type Validator interface {
	ValidateTokenWithClaims(token string, claims map[string]string) error
}

func NewValidator(jwtVerificationKey string) (Validator, error) {
	key, err := parseJWTVerificationKey(jwtVerificationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %s", err.Error())
	}
	return &validator{jwtVerificationKey: key}, nil
}

func (s *validator) ValidateTokenWithClaims(tokenString string, claims map[string]string) error {
	var tokenClaims jwt.MapClaims
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if _, signedUsingRSA := token.Method.(*jwt.SigningMethodRSA); !signedUsingRSA {
			return nil, errors.New("token not signed using RSA")
		}

		return s.jwtVerificationKey, nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse token: %s", err.Error())
	}

	if !token.Valid {
		return errors.New("token is invalid")
	}

	if !strings.HasPrefix(tokenClaims["iss"].(string), claims["iss"]) {
		return errors.New("iss claim does not match")
	}

	return nil
}

func parseJWTVerificationKey(jwtVerificationKey string) (*rsa.PublicKey, error) {
	key := strings.TrimSpace(jwtVerificationKey)
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
