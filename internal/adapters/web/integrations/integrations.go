package integrations

import (
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/sso"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type Integrations struct {
	SSO sso.Integration
}

func CreateIntegration(cfg *config.ConfigurationService) *Integrations {
	return &Integrations{
		SSO: sso.NewIntegration(cfg),
	}
}
