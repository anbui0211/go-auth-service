package services

import (
	"context"

	"goauth/internal/dao"
	"goauth/pkg/cache"
	"time"

	"gorm.io/gorm"
)

type IUserService interface {
	ChangeStatus(db *gorm.DB, userId string) error
}

type userService struct {
	userDao dao.IUserDao
}

func NewuserService(userDao dao.IUserDao) IUserService {
	return userService{
		userDao: userDao,
	}
}

// ChangeStatus implements IUserService.
func (us userService) ChangeStatus(db *gorm.DB, userId string) error {
	newStatus, err := us.userDao.UpdateStatus(db, userId)
	if err != nil {
		return err
	}

	var (
		// "user_status" + userID
		keyUserStatusRedis = cache.GenKeyRedis("user_status", userId)
		timeExpired        = time.Hour * 24 * 7
	)

	if err := cache.SetRedis(context.Background(), keyUserStatusRedis, newStatus, timeExpired); err != nil {
		return err
	}

	return nil
}
