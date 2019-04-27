package handler

const (
	ErrInvalidToken                = "invalid_token"
	ErrDescriptionNoAuthentication = "Full authentication is required to access this resource"
	ErrDescriptionMalformedToken   = "The request token is malformed. Please validate that your request token was issued by the UAA server authorized by CredHub."
)
