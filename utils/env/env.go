package env

import "os"

func EnvOauth2RedirectURL() string {
	return os.Getenv("OAUTH2_REDIRECT_URL")
}
func EnvOauth2ClientID() string {
	return os.Getenv("OAUTH2_CLIENT")
}
func EnvOauth2ClientSecret() string {
	return os.Getenv("OAUTH2_CLIENT_SECRET")
}
