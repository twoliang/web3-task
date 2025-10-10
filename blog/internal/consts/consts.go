package consts

import (
	"golang.org/x/crypto/bcrypt"
)

const JWTSecret = "k8X9n7Z6m5B4a3S2d1F0g9H8j7K6l5P4"

// HashPassword 对密码进行加密并返回哈希值
func HashPassword(password string) (string, error) {
	// 生成密码哈希，第二个参数是成本因子，值越大加密越慢但安全性越高
	// 通常建议使用 10-14 之间的值
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// CheckPassword 验证密码与哈希值是否匹配
func CheckPassword(password, hash string) error {
	// 比较密码与哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
