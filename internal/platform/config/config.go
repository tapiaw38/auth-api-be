package config

var configService *ConfigurationService

const (
	ReleaseMode GinModeServer = "release"
	DebugMode   GinModeServer = "debug"
)

type (
	GinModeServer string

	ConfigurationService struct {
		ServerConfig ServerConfig
		DBConfig     DBConfig
		S3Config     S3Config
		GCPConfig    GCPConfig
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
)

func InitConfigService(config *ConfigurationService) {
	if configService == nil {
		configService = config
	}
}

func GetConfigService() ConfigurationService {
	return *configService
}
