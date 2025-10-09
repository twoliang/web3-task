package logic

import (
	"blog/internal/model"
	"errors"

	"gorm.io/gorm"
)

// service/user.go
type UserService struct {
	DB *gorm.DB
}

func (s *UserService) Register(username, email, password string) (*model.User, error) {
	// 检查用户名/邮箱是否已存在
	var existingUser model.User
	if err := s.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}
	if err := s.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 创建用户
	newUser := model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.DB.Create(&newUser).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return &newUser, nil
}
