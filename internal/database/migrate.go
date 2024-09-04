package database

import (
	"fmt"

	authjwt "goauth/internal/auth/jwt"
	constants "goauth/internal/constant"
	gormmodel "goauth/internal/models/gorm"
	urand "goauth/pkg/utils/rand"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&gormmodel.User{})

	user := gormmodel.User{
		UserID:    urand.RandUuid(),
		FirstName: "Bui",
		LastName:  "An",
		Email:     "anbui@gmail.com",
		Role:      "ADMIN",
		Status:    constants.StatusActive,
	}

	hashPass, err := authjwt.HassPassword("1234")
	if err != nil {
		fmt.Println("Error hashing password")
	}
	user.Password = hashPass

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("Error seeding data:", err)
	}

	fmt.Println("Auto Migration has beed processed")
}
