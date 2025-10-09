package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论信息表
type Comment struct {
	gorm.Model
	ID        uint       `gorm:"primary_key;auto_increment;comment:评论ID" json:"id"`
	Content   string     `gorm:"type:text;not_null;comment:评论内容" json:"content"`
	UserID    uint       `gorm:"type:int;not_null;comment:用户ID" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID;references:ID;comment:评论用户"` // 外键关联
	PostID    uint       `gorm:"type:int;not_null;comment:文章ID" json:"post_id"`
	Post      Post       `gorm:"foreignKey:PostID;references:ID;comment:评论文章"` // 外键关联
	CreatedAt time.Time  `gorm:"type:timestamp;not_null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp;not_null;default:CURRENT_TIMESTAMP;on_update:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp;default:null;comment:删除时间" json:"deleted_at"`
}
