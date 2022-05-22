package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorId      uint
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
}
