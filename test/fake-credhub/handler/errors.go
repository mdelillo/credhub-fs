package handler

const (
	ErrCredentialDoesNotExist      = "The request could not be completed because the credential does not exist or you do not have sufficient authorization."
	ErrDescriptionMalformedToken   = "The request token is malformed. Please validate that your request token was issued by the UAA server authorized by CredHub."
	ErrDescriptionNoAuthentication = "Full authentication is required to access this resource"
	ErrInvalidPathOrBody           = "The request could not be fulfilled because the request path or body did not meet expectation. Please check the documentation for required formatting and retry your request."
	ErrInvalidToken                = "invalid_token"
	ErrInvalidType                 = "Only 'value' types are supported"
	ErrMissingNameParameter        = "The query parameter name is required for this request."
)
