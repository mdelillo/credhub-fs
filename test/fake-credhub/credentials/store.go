package credentials

import "strings"

type store struct {
	credentials map[string]Credential
}

type Store interface {
	GetByName(name string) (cred Credential, found bool)
	GetByPath(path string) []Credential
	Set(credential Credential)
}

func NewStore() Store {
	return &store{
		credentials: map[string]Credential{},
	}
}

func (s *store) GetByName(name string) (Credential, bool) {
	for n, cred := range s.credentials {
		if name == n {
			return cred, true
		}
	}
	return Credential{}, false
}

func (s *store) GetByPath(path string) []Credential {
	var matchingCredentials []Credential

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	for name, cred := range s.credentials {
		if strings.HasPrefix(name, path) {
			matchingCredentials = append(matchingCredentials, cred)
		}
	}

	return matchingCredentials
}

func (s *store) Set(credential Credential) {
	s.credentials[credential.Name] = credential
}
