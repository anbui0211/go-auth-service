package dao

import (
	"log"

	gormmodel "goauth/internal/models/gorm"

	"gorm.io/gorm"
)

type IUserDao interface {
	FindByEmail(db *gorm.DB, email string) (user gormmodel.User, err error)
	FindByID(db *gorm.DB, id string) (user gormmodel.User, err error)
	CountByEmail(db *gorm.DB, email string) (count int64)
	Create(db *gorm.DB, user gormmodel.User) (gormmodel.User, error)
}

type userDao struct{}

func NewUserDao() IUserDao {
	return &userDao{}
}

func (d *userDao) FindByEmail(db *gorm.DB, email string) (user gormmodel.User, err error) {
	err = db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println("error finding user: ", err)
		return gormmodel.User{}, err
	}

	return
}

func (d *userDao) FindByID(db *gorm.DB, id string) (user gormmodel.User, err error) {
	err = db.Where("user_id = ?", id).First(&user).Error
	if err != nil {
		log.Println("error finding user: ", err)
		return gormmodel.User{}, err
	}
	return
}

func (d *userDao) CountByEmail(db *gorm.DB, email string) (count int64) {
	db.Model(&gormmodel.User{}).Where("email = ?", email).Count(&count)
	return
}

func (d *userDao) Create(db *gorm.DB, user gormmodel.User) (gormmodel.User, error) {
	if err := db.Create(&user).Error; err != nil {
		return gormmodel.User{}, err
	}

	return user, nil
}
