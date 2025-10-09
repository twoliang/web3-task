package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户信息表
type User struct {
	gorm.Model
	ID        uint       `gorm:"primary_key;auto_increment;comment:用户ID" json:"id"`
	Username  string     `gorm:"type:varchar(50);not_null;unique;comment:用户名" json:"username"`
	Password  string     `gorm:"type:varchar(255);not_null;comment:密码" json:"password"`
	Email     string     `gorm:"type:varchar(100);not_null;unique;comment:电子邮箱" json:"email"`
	CreatedAt time.Time  `gorm:"type:timestamp;not_null;default:CURRENT_TIMESTAMP;comment:记录创建时间" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:timestamp;not_null;default:CURRENT_TIMESTAMP;on_update:CURRENT_TIMESTAMP;comment:记录更新时间" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp;default:null;comment:删除时间" json:"deleted_at"`
}
