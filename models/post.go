package models

import "gorm.io/gorm"

type PostBase struct {
}

type Post struct {
	gorm.Model
	*PostBase
	Title   string `gorm:"size:255;not null"`
	Content string `gorm:"type:varchar;not null"`
	UserID  uint   `gorm:"index"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tags    []*Tag `gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
