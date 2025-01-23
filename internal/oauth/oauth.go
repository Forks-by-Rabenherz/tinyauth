package oauth

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func NewOAuth(config oauth2.Config) *OAuth {
	return &OAuth{
		Config: config,
	}
}

type OAuth struct {
	Config   oauth2.Config
	Context  context.Context
	Token    *oauth2.Token
	Verifier string
}

func (oauth *OAuth) Init() {
	oauth.Context = context.Background()
	oauth.Verifier = oauth2.GenerateVerifier()
}

func (oauth *OAuth) GetAuthURL() string {
	return oauth.Config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(oauth.Verifier))
}

func (oauth *OAuth) ExchangeToken(code string) error {
	token, err := oauth.Config.Exchange(oauth.Context, code, oauth2.VerifierOption(oauth.Verifier))
	if err != nil {
		log.Error().Err(err).Msg("Failed to exchange code")
		return err
	}
	oauth.Token = token
	return nil
}

func (oauth *OAuth) GetClient() *http.Client {
	return oauth.Config.Client(oauth.Context, oauth.Token)
}
