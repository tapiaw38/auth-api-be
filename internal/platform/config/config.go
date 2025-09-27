package config

var configService *ConfigurationService

const (
	ReleaseMode GinModeServer = "release"
	DebugMode   GinModeServer = "debug"
)

type (
	GinModeServer string

	ConfigurationService struct {
		AppName      string
		ServerConfig ServerConfig
		DBConfig     DBConfig
		S3Config     S3Config
		GCPConfig    GCPConfig
		RabbitMQ     RabbitMQConfig
		Notification NotificationConfig
		InitConfig   InitConfig
	}

	ServerConfig struct {
		GinMode   GinModeServer
		Port      string
		Host      string
		JWTSecret string
	}

	DBConfig struct {
		DatabaseURL string
	}

	S3Config struct {
		AWSRegion          string
		AWSAccessKeyID     string
		AWSSecretAccessKey string
		AWSBucket          string
	}

	OAuth2Config struct {
		GoogleClientID     string
		GoogleClientSecret string
		FrontendURL        string
	}

	GCPConfig struct {
		OAuth2Config OAuth2Config
	}

	RabbitMQConfig struct {
		URL string
	}

	EmailConfig struct {
		Host     string
		Port     string
		Username string
		Password string
	}

	NotificationConfig struct {
		Email EmailConfig
	}
	
	InitConfig struct {
		EnsureDefaultRoles bool
	}
)

func InitConfigService(config *ConfigurationService) {
	if configService == nil {
		configService = config
	}
}

func GetConfigService() ConfigurationService {
	return *configService
}
