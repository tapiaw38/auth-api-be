package appcontext

import (
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources"
	"github.com/tapiaw38/auth-api-be/internal/adapters/datasources/repositories"
	"github.com/tapiaw38/auth-api-be/internal/adapters/queue"
	"github.com/tapiaw38/auth-api-be/internal/adapters/web/integrations"
	"github.com/tapiaw38/auth-api-be/internal/platform/config"
)

type Context struct {
	Repositories  *repositories.Repositories
	Integrations  *integrations.Integrations
	RabbitMQ      *queue.RabbitMQ
	ConfigService *config.ConfigurationService
}

type Option func(*Context)

type Factory func(opts ...Option) *Context

func NewFactory(
	datasources *datasources.Datasources,
	integrations *integrations.Integrations,
	rabbitMQ *queue.RabbitMQ,
	configService *config.ConfigurationService,
) func(opts ...Option) *Context {
	return func(opts ...Option) *Context {
		return &Context{
			Repositories:  repositories.NewFactory(datasources, configService)(),
			Integrations:  integrations,
			RabbitMQ:      rabbitMQ,
			ConfigService: configService,
		}
	}
}
