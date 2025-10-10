package logic

import (
	"errors"
	"fmt"
	"time"
	"web-task/blog/internal/consts"
	"web-task/blog/internal/model"

	"github.com/golang-jwt/jwt"
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
	hashedPassword, err := consts.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 创建用户
	fmt.Printf("Start to register with username: %s, password: %s, email: %s\n", username, password, email)
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

func (s *UserService) Login(username, password string) (token string, err error) {
	// 1. 检查用户是否存在
	var user model.User
	if err = s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	// 2. 验证密码
	if err = consts.CheckPassword(password, user.Password); err != nil {
		return "", errors.New("invalid password")
	}

	// 3. 生成 JWT 令牌
	token, err = s.generateJWT(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *UserService) generateJWT(user model.User) (string, error) {
	// 设置 JWT 的有效载荷
	claims := jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 设置令牌有效期为 24 小时
	}

	// 使用 HS256 算法签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并返回令牌
	signedToken, err := token.SignedString([]byte(consts.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
