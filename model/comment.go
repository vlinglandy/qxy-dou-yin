package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserId     uint
	Content    string
	CreateDate string
	VideoId    uint
}
