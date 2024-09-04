package authjwt

import (
	"errors"
	"fmt"
	"goauth/internal/auth"
	urand "goauth/pkg/utils/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access_token"
	TokenTypeRefresh = "refresh_token"
)

func CreateAccessToken(user auth.User) (string, error) {
	// Create a new JWT with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"role": user.Role,
		"iss":  "auth-app",                       // Issuer
		"iat":  time.Now().Unix(),                // Issued at
		"exp":  time.Now().Add(time.Hour).Unix(), // Expiration time
	})

	privateKey, err := getPrivateKey()
	if err != nil {
		return "", err
	}
	tokenString, err := claims.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(user auth.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": user.ID,
		"iss": "auth-app",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expirations
	})

	privateKey, err := getPrivateKey()
	if err != nil {
		return "", err
	}

	tokenString, err := claims.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// validation signing method
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// Get public key and return it
		publicKey, err := GetPublicKey()
		if err != nil {
			return nil, fmt.Errorf("get public key fail %v", err)
		}

		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check valid claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" || !urand.IsValidUuid(userID) {
			return nil, errors.New("invalid token: missing user ID")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, errors.New("token expried")
			}
		}
	}

	return token, nil
}

func VerifyTokenV2(tokenString string, tokenType string) (userId, userName, userRole string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// validation signing method
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		// Get public key and return it
		publicKey, err := GetPublicKey()
		if err != nil {
			return nil, fmt.Errorf("get public key fail %v", err)
		}

		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return "", "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", "", errors.New("fail to covert token claim to jwt.MapClaims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return "", "", "", errors.New("token expired")
		}
	}

	userId, ok = claims["sub"].(string)
	if !ok || userId == "" {
		return "", "", "", errors.New("user id does not exist in token")
	}

	// No handle if is refresh token
	if tokenType == TokenTypeAccess {
		userName, ok = claims["name"].(string)
		if !ok || userName == "" {
			return "", "", "", errors.New("user name does not exist in token")
		}

		userRole, ok = claims["role"].(string)
		if !ok || userRole == "" {
			return "", "", "", errors.New("user role does not exist in token")
		}
	}

	return userId, userName, userRole, nil
}
