package auth

import (
	"fmt"

	witAuth "github.com/hayride-dev/bindings/go/internal/gen/imports/hayride/mcp/auth"

	"go.bytecodealliance.org/cm"
)

var _ Provider = (*ProviderResource)(nil)

type Provider interface {
	AuthURL() (string, error)
	Registration(data []byte) ([]byte, error)
	ExchangeCode(data []byte) ([]byte, error)
	Validate(token string) (bool, error)
}

type ProviderResource cm.Resource

func New() (ProviderResource, error) {
	return ProviderResource(witAuth.NewProvider()), nil
}

func (t ProviderResource) AuthURL() (string, error) {
	provider := cm.Reinterpret[witAuth.Provider](t)

	result := provider.AuthURL()
	if result.IsErr() {
		return "", fmt.Errorf("failed to get auth URL: %s", result.Err().Data())
	}

	return *result.OK(), nil
}

func (t ProviderResource) Registration(data []byte) ([]byte, error) {
	provider := cm.Reinterpret[witAuth.Provider](t)

	result := provider.Registration(cm.ToList(data))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to get registration URL: %s", result.Err().Data())
	}

	return result.OK().Slice(), nil
}

func (t ProviderResource) ExchangeCode(data []byte) ([]byte, error) {
	provider := cm.Reinterpret[witAuth.Provider](t)

	result := provider.ExchangeCode(cm.ToList(data))
	if result.IsErr() {
		return nil, fmt.Errorf("failed to exchange code: %s", result.Err().Data())
	}

	return result.OK().Slice(), nil
}

func (t ProviderResource) Validate(token string) (bool, error) {
	provider := cm.Reinterpret[witAuth.Provider](t)

	result := provider.Validate(token)
	if result.IsErr() {
		return false, fmt.Errorf("failed to validate token: %s", result.Err().Data())
	}

	return *result.OK(), nil
}
