package models

import "gorm.io/gorm"

type TagBase struct {
}

type Tag struct {
	gorm.Model
	*TagBase
	Name   string  `gorm:"size:255;not null"`
	Posts  []*Post `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User   User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID uint    `gorm:"index"`
}
