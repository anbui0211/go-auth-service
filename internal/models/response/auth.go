package responsemodel

import "goauth/internal/auth"

type ResponseToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResponseAuth struct {
	Token ResponseToken `json:"token"`
	User  auth.User     `json:"user"`
}
