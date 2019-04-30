package util

import (
	"github.com/mdelillo/credhub-fs/pkg/credhub"
)

type Dependencies interface {
	GetCredhubClient() credhub.Client
	SetCredhubClient(credhub.Client)
}

type dependencies struct {
	credhubClient credhub.Client
}

func NewDependencies() Dependencies {
	return &dependencies{}
}

func (c *dependencies) SetCredhubClient(credhubClient credhub.Client) {
	c.credhubClient = credhubClient
}

func (c *dependencies) GetCredhubClient() credhub.Client {
	return c.credhubClient
}
