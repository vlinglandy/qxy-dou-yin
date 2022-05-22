package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username      string
	Password      string
	FollowCount   int64
	FollowerCount int64
}
