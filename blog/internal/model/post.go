package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID        uint       `gorm:"primary_key;auto_increment;comment:文章ID" json:"id"`
	Title     string     `gorm:"type:varchar(200);not_null;comment:文章标题" json:"title"`
	Content   string     `gorm:"type:text;not_null;comment:文章内容" json:"content"`
	UserID    uint       `gorm:"type:int;not_null;" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID;references:ID" json:"user"` // 添加用户关系
	CreatedAt time.Time  `gorm:"type:timestamp;not_null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp;not_null;default:CURRENT_TIMESTAMP;on_update:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp;default:null;comment:删除时间" json:"deleted_at"`
}
