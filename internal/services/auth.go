package services

import (
	"errors"
	"fmt"

	constants "goauth/internal/constant"
	"goauth/internal/dao"
	gormmodel "goauth/internal/models/gorm"
	requestmodel "goauth/internal/models/request"
	responsemodel "goauth/internal/models/response"
	uauth "goauth/utils/auth"
	ujwt "goauth/utils/auth/jwt"
	urand "goauth/utils/rand"

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
		Role:      constants.RoleAdmin,
		Status:    constants.StatusActive,
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

	// Check if user is banned
	if user.Status == constants.StatusInActive {
		return responsemodel.ResponseAuth{}, errors.New("user is banned, can not login")
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
	userId, _, _, err := ujwt.VerifyTokenV2(payload.RefreshToken, "refresh_token")
	if err != nil {
		return "", err
	}

	user, err := as.userDao.FindByID(db, userId)
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
