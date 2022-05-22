package model

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	UserId   uint
	ToUserId uint
}
