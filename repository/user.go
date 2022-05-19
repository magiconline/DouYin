package repository

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// UserId int64
	Email         string `gorm:column"type:varchar(32);unique;not null;index"`
	UserName      string `gorm:"type:varchar(32);unique;not null"`
	Password      string `gorm:"type:varchar(32)"`
	Token         string
	FollowCount   int `gorm:"default:0"`
	FollowerCount int `gorm:"default:0"`
}
