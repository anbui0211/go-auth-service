package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"goauth/internal/auth"
	authjwt "goauth/internal/auth/jwt"
	constants "goauth/internal/constant"
	"goauth/internal/dao"
	gormmodel "goauth/internal/models/gorm"
	requestmodel "goauth/internal/models/request"
	responsemodel "goauth/internal/models/response"
	"goauth/pkg/cache"
	urand "goauth/pkg/utils/rand"

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

	hashPassword, err := authjwt.HassPassword(payload.Password)
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

	// Set status user to redis
	var (
		// "user_status" + userID
		keyUserStatusRedis = cache.GenKeyRedis("user_status", userCreate.UserID)
		timeExpired        = time.Hour * 24 * 7
	)
	if err := cache.SetRedis(context.Background(), keyUserStatusRedis, userCreate.Status, timeExpired); err != nil {
		return err
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
	if ok := authjwt.VerifyPassword(payload.Password, user.Password); !ok {
		return responsemodel.ResponseAuth{}, errors.New("password does not match")
	}

	// create a token
	authUser := auth.User{
		ID:   user.UserID,
		Name: fmt.Sprintf("%s %s", user.LastName, user.FirstName),
		Role: user.Role,
	}

	accessToken, err := authjwt.CreateAccessToken(authUser)
	if err != nil {
		return responsemodel.ResponseAuth{}, errors.New("create access token failed")
	}

	refreshToken, err := authjwt.CreateRefreshToken(authUser)
	if err != nil {
		return responsemodel.ResponseAuth{}, errors.New("create refresh token failed")
	}

	// Set status user to redis
	var (
		// "user_status" + userID
		keyUserStatusRedis = cache.GenKeyRedis("user_status", user.UserID)
		timeExpired        = time.Hour * 24 * 7
	)

	if err := cache.SetRedis(context.Background(), keyUserStatusRedis, user.Status, timeExpired); err != nil {
		return responsemodel.ResponseAuth{}, err
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
	userId, _, _, err := authjwt.VerifyTokenV2(payload.RefreshToken, "refresh_token")
	if err != nil {
		return "", err
	}

	user, err := as.userDao.FindByID(db, userId)
	if err != nil {
		return "", errors.New("refresh token: user ID invalid")
	}

	accessToken, err := authjwt.CreateAccessToken(auth.User{
		ID:   user.UserID,
		Name: fmt.Sprintf("%s %s", user.LastName, user.FirstName),
		Role: user.Role,
	})
	if err != nil {
		return "", errors.New("create token failed")
	}

	return accessToken, nil
}
