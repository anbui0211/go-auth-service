package gormmodel

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID    string `gorm:"column:user_id"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	Role      string `gorm:"column:role"`
	Status    string `gorm:"column:status"`
}
