package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email          string `gorm:"size:255;not null" validate:"email"`
	PasswordDigest string `gorm:"size:255;not null" json:"-"`
	Posts          []Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags           []Tag  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserLogin struct {
	Email    string `validate:"email"`
	Password string
}
