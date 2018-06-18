package credhub_shims

import (
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

//go:generate counterfeiter -o ./credhub_fakes/credhub_auth_fake.go . CredhubAuth
type CredhubAuth interface {
	UaaClientCredentials(clientId, clientSecret string) auth.Builder
}

type CredhubAuthShim struct {
}

func (c *CredhubAuthShim) UaaClientCredentials(clientId, clientSecret string) auth.Builder {
	return auth.UaaClientCredentials(clientId, clientSecret)
}

//go:generate counterfeiter -o ./credhub_fakes/credhub_fake.go . Credhub
type Credhub interface {
	SetJSON(name string, value values.JSON, overwrite credhub.Mode) (credentials.JSON, error)
	GetLatestJSON(name string) (credentials.JSON, error)
	Delete(name string) error
}

type CredhubShim struct {
	delegate *credhub.CredHub
}

func NewCredhubShim(
	url string,
	caCert string,
	clientID string,
	clientSecret string,
	authShim CredhubAuth,
) (Credhub, error) {
	delegate, err := credhub.New(
		url,
		credhub.CaCerts(caCert),
		credhub.SkipTLSValidation(false),
		credhub.Auth(authShim.UaaClientCredentials(clientID, clientSecret)),
	)
	if err != nil {
		return nil, err
	}

	return &CredhubShim{
		delegate: delegate,
	}, nil
}

func (ch *CredhubShim) SetJSON(name string, value values.JSON, overwrite credhub.Mode) (credentials.JSON, error) {
	return ch.delegate.SetJSON(name, value, overwrite)
}

func (ch *CredhubShim) GetLatestJSON(name string) (credentials.JSON, error) {
	return ch.delegate.GetLatestJSON(name)
}

func (ch *CredhubShim) Delete(name string) error {
	return ch.delegate.Delete(name)
}
