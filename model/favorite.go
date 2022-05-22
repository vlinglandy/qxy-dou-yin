package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserId  uint
	VideoId uint
}
