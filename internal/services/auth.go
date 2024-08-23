package services

import (
	"errors"
	"fmt"

	"goauth/internal/dao"
	gormmodel "goauth/internal/models/gorm"
	requestmodel "goauth/internal/models/request"
	responsemodel "goauth/internal/models/response"
	uauth "goauth/utils/auth"
	ujwt "goauth/utils/auth/jwt"
	urand "goauth/utils/rand"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type IAuthService interface {
	Register(db *gorm.DB, payload requestmodel.RegisterPayload) error
	Login(db *gorm.DB, payload requestmodel.LoginPayload) (responsemodel.ResponseAuth, error)
	RefreshToken(db *gorm.DB, payload requestmodel.RefreshPayload) (string, error)
	HandleGoogleOAuthCallback(db *gorm.DB, state, code string) (*responsemodel.ResponseAuth, error)
}
type authService struct {
	userDao dao.IUserDao
}

func NewAuthService(userDao dao.IUserDao) IAuthService {
	return &authService{
		userDao: userDao,
	}
}

func (as *authService) Register(db *gorm.DB, payload requestmodel.RegisterPayload) error {
	// check if user already registered
	count := as.userDao.CountByEmail(db, payload.Email)
	if count > 0 {
		return errors.New("email already registered")
	}

	hashPassword, err := ujwt.HassPassword(payload.Password)
	if err != nil {
		return errors.New("hash password error: ")
	}

	userCreate := gormmodel.User{
		UserID:    urand.RandUuid(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashPassword,
		Role:      "USER",
	}

	if _, err := as.userDao.Create(db, userCreate); err != nil {
		return errors.New("error creating user")
	}

	return nil
}

func (as *authService) Login(db *gorm.DB, payload requestmodel.LoginPayload) (responsemodel.ResponseAuth, error) {
	// Check if user does not exist
	user, err := as.userDao.FindByEmail(db, payload.Email)
	if err != nil {
		return responsemodel.ResponseAuth{}, errors.New("email does not exist")
	}

	// Verify password
	if ok := ujwt.VerifyPassword(payload.Password, user.Password); !ok {
		return responsemodel.ResponseAuth{}, errors.New("password does not match")
	}

	// create a token
	authUser := uauth.User{
		ID:   user.UserID,
		Name: fmt.Sprintf("%s %s", user.LastName, user.FirstName),
		Role: user.Role,
	}

	accessToken, err := ujwt.CreateAccessToken(authUser)
	if err != nil {
		return responsemodel.ResponseAuth{}, errors.New("create access token failed")
	}

	refreshToken, err := ujwt.CreateRefreshToken(authUser)
	if err != nil {
		return responsemodel.ResponseAuth{}, errors.New("create refresh token failed")
	}

	return responsemodel.ResponseAuth{
		Token: responsemodel.ResponseToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		User: authUser,
	}, nil
}

func (as *authService) RefreshToken(db *gorm.DB, payload requestmodel.RefreshPayload) (string, error) {
	token, err := ujwt.VerifyToken(payload.RefreshToken)
	if err != nil {
		return "", err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	userID, _ := claims["sub"].(string)
	user, err := as.userDao.FindByID(db, userID)
	if err != nil {
		return "", errors.New("refresh token: user ID invalid")
	}

	accessToken, err := ujwt.CreateAccessToken(uauth.User{
		ID:   user.UserID,
		Name: fmt.Sprintf("%s %s", user.LastName, user.FirstName),
		Role: user.Role,
	})
	if err != nil {
		return "", errors.New("create token failed")
	}

	return accessToken, nil
}
