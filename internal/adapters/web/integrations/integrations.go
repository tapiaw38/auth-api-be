package integrations

import (
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/notification"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations/sso"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type Integrations struct {
	SSO          sso.Integration
	Notification notification.Integration
}

func CreateIntegration(cfg *config.ConfigurationService) *Integrations {
	return &Integrations{
		SSO:          sso.NewIntegration(cfg),
		Notification: notification.NewIntegration(cfg),
	}
}
