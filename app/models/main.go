package models

import (
	"gorm.io/gorm"
)

// type Post struct {
// 	gorm.Model
// 	Title   string `gorm:"size:255;not null"`
// 	Content string `gorm:"type:varchar;not null"`
// }

// type Tag struct {
// 	gorm.Model
// 	Name string `gorm:"size:255;not null"`
// }

type userBase struct {
	Email string `validate:"email"`
}

type User struct {
	gorm.Model
	Email          string `gorm:"size:255;not null" validate:"email"`
	PasswordDigest string `gorm:"size:255;not null" json:"-"`
}

type UserLogin struct {
	userBase
	Password string
}
