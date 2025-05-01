package sso

import (
	"context"

	"github.com/tapiaw38/auth-api-be/internal/platform/config"
	"golang.org/x/oauth2"
	googleauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type (
	Integration interface {
		ExchangeCode(context.Context, string) (*oauth2.Token, error)
		GetUserInfo(context.Context, *oauth2.Token) (*SocialUser, error)
	}

	integration struct {
		config *oauth2.Config
	}

	SocialUser struct {
		Token         string `json:"token"`
		RefreshToken  string `json:"refresh_token"`
		Scopes        string `json:"scopes"`
		Email         string `json:"email"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		Picture       string `json:"picture"`
		Birthday      string `json:"birthday"`
		VerifiedEmail bool   `json:"verified_email"`
	}
)

func NewIntegration(cfg *config.ConfigurationService) Integration {
	config := initConfig(cfg)
	return &integration{config: config}
}

func initConfig(cfg *config.ConfigurationService) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.GCPConfig.OAuth2Config.GoogleClientID,
		ClientSecret: cfg.GCPConfig.OAuth2Config.GoogleClientSecret,
		RedirectURL:  cfg.GCPConfig.OAuth2Config.FrontendURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
}

func (i *integration) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := i.config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (*integration) GetUserInfo(ctx context.Context, token *oauth2.Token) (*SocialUser, error) {
	service, err := googleauth.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(token)))
	if err != nil {
		return nil, err
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	return &SocialUser{
		Email:         userInfo.Email,
		FirstName:     userInfo.GivenName,
		LastName:      userInfo.FamilyName,
		Picture:       userInfo.Picture,
		VerifiedEmail: true,
	}, nil
}
