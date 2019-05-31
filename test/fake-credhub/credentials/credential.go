package credentials

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Type             string    `json:"type"`
	Value            string    `json:"value"`
	VersionCreatedAt time.Time `json:"version_created_at"`
}

type CredentialNameAndDate struct {
	Name             string    `json:"name"`
	VersionCreatedAt time.Time `json:"version_created_at"`
}
