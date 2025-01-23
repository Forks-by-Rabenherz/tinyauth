package providers

import (
	"tinyauth/internal/oauth"
	"tinyauth/internal/types"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func NewProviders(config types.OAuthConfig) *Providers {
	return &Providers{
		Config: config,
	}
}

type Providers struct {
	Config    types.OAuthConfig
	Github    *oauth.OAuth
	Google    *oauth.OAuth
	Microsoft *oauth.OAuth
}

func (providers *Providers) Init() {
	if providers.Config.GithubClientId != "" && providers.Config.GithubClientSecret != "" {
		log.Info().Msg("Initializing Github OAuth")
		providers.Github = oauth.NewOAuth(oauth2.Config{
			ClientID:     providers.Config.GithubClientId,
			ClientSecret: providers.Config.GithubClientSecret,
			Scopes:       GithubScopes(),
			Endpoint:     endpoints.GitHub,
		})
		providers.Github.Init()
	}
}

func (providers *Providers) Login(code string, provider string) (string, error) {
	switch provider {
	case "github":
		if providers.Github == nil {
			return "", nil
		}
		exchangeErr := providers.Github.ExchangeToken(code)
		if exchangeErr != nil {
			return "", exchangeErr
		}
		client := providers.Github.GetClient()
		email, emailErr := GetGithubEmail(client)
		if emailErr != nil {
			return "", emailErr
		}
		return email, nil
	default:
		return "", nil
	}
}

func (providers *Providers) GetUser(provider string) (string, error) {
	switch provider {
	case "github":
		if providers.Github == nil {
			return "", nil
		}
		client := providers.Github.GetClient()
		email, emailErr := GetGithubEmail(client)
		if emailErr != nil {
			return "", emailErr
		}
		return email, nil
	default:
		return "", nil
	}
}

func (providers *Providers) GetAuthURL(provider string) string {
	switch provider {
	case "github":
		if providers.Github == nil {
			return ""
		}
		return providers.Github.GetAuthURL()
	default:
		return ""
	}
}
