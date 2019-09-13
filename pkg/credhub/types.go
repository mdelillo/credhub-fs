package credhub

import (
	"time"

	"github.com/google/uuid"
)

// TODO: Can `Value` be interface{} and handle arbitrary JSON?
type Credential struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Value            string    `json:"value"`
	VersionCreatedAt time.Time `json:"version_created_at"`
}

type ErrCredentialNotFound struct {
	credentialName string
}

func (e *ErrCredentialNotFound) Error() string {
	return "could not find credential " + e.credentialName
}
