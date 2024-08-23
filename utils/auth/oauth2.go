package uauth

import (
	"goauth/utils/env"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OauthStateString = "random-string" // Please generate different string

func GetGoogleOauthConfig() oauth2.Config {
	return oauth2.Config{
		RedirectURL:  env.EnvOauth2RedirectURL(),
		ClientID:     env.EnvOauth2ClientID(),
		ClientSecret: env.EnvOauth2ClientSecret(),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	}
}
