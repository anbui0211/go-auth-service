package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	gormmodel "goauth/internal/models/gorm"
	responsemodel "goauth/internal/models/response"
	uauth "goauth/utils/auth"
	ujwt "goauth/utils/auth/jwt"
	urand "goauth/utils/rand"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

func (as *authService) HandleGoogleOAuthCallback(db *gorm.DB, state, code string) (*responsemodel.ResponseAuth, error) {
	//  Get Google OAuth configuration
	config := uauth.GetGoogleOauthConfig()

	// Validate OAuth state
	if state != uauth.OauthStateString {
		log.Printf("Invalid oauth state, expected '%s', got '%s'\n", uauth.OauthStateString, state)
		return nil, errors.New("invalid oauth state")
	}

	// Exchange the authorization code for a token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("code exchange failed: %s\n", err.Error())
		return nil, errors.New("code exchange failed")
	}

	// Extract the ID token from the OAuth token
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("failed to extrasct ID token")
	}

	// Validate the ID token
	payload, err := idtoken.Validate(context.Background(), idToken, config.ClientID)
	if err != nil {
		log.Printf("Token validation failed: %s\n", err.Error())
		return nil, errors.New("token validation failed")
	}

	// save if user have not been registered
	var userInfoRes uauth.User

	user, err := as.userDao.FindByEmail(db, payload.Claims["email"].(string))
	if err != nil {
		// If user has not been registered
		userCreate := gormmodel.User{
			UserID:    urand.RandUuid(),
			FirstName: payload.Claims["name"].(string),
			Email:     payload.Claims["email"].(string),
			Role:      "USER",
		}

		userCreate, err := as.userDao.Create(db, userCreate)
		if err != nil {
			return nil, errors.New("error creating user")
		}

		userInfoRes = uauth.User{
			ID:   userCreate.UserID,
			Name: fmt.Sprintf("%s %s", userCreate.LastName, userCreate.FirstName),
			Role: userCreate.Role,
		}
	} else {
		// If user has been registered
		userInfoRes = uauth.User{
			ID:   user.UserID,
			Name: fmt.Sprintf("%s %s", user.LastName, user.FirstName),
			Role: user.Role,
		}
	}

	accessToken, err := ujwt.CreateAccessToken(userInfoRes)
	if err != nil {
		return nil, errors.New("create access token failed")
	}

	refreshToken, err := ujwt.CreateRefreshToken(userInfoRes)
	if err != nil {
		return nil, errors.New("create refresh token failed")

	}

	res := &responsemodel.ResponseAuth{
		Token: responsemodel.ResponseToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: userInfoRes,
	}

	return res, nil
}
