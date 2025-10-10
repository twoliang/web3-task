package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"errors"
	"web-task/blog/internal/consts"
	"web-task/blog/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// AuthMiddleware 是一个中间件函数，用于验证 JWT 令牌
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 解析 Bearer 令牌
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 解析 JWT 令牌
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 确保令牌的签名方法是 HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(consts.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 检查令牌是否有效
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 从 Claims 中获取用户名
			username := claims["username"].(string)

			// 将用户名存储到上下文中，供后续使用
			c.Set("username", username)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
		}
	}
}

// CheckBlogOwnerMiddleware 是一个中间件函数，用于检查请求者是否是博客的主人
func CheckBlogOwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			blog model.Post
			u    model.User
			DB   *gorm.DB
		)

		// 从上下文中获取用户名
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// 从请求路径中获取博客 ID（确保是合法正整数）
		blogIDStr := c.Param("id")
		if blogIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID is required"})
			c.Abort()
			return
		}
		blogID, err := strconv.ParseUint(blogIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Blog ID"})
			c.Abort()
			return
		}

		// 开始事务（可选，保证“存在性+归属”一致性）
		DB = c.MustGet("db").(*gorm.DB)
		tx := DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// 1) 只取主键 ID，减少字段扫描
		if err := tx.Select("id").Where("username = ?", username).First(&u).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			} else {
				// 可选：记录错误日志
				// log.Printf("check blog owner: query user failed: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			c.Abort()
			return
		}

		// 2) 用用户 ID 与博客 ID 一次性校验归属（避免先查博客再比对的并发问题）
		if err := tx.Where("id = ? AND user_id = ?", blogID, u.ID).First(&blog).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found or you are not the owner"})
			} else {
				// 可选：记录错误日志
				// log.Printf("check blog owner: query blog failed: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
			c.Abort()
			return
		}

		// 可选：将当前用户 ID 注入上下文，便于后续处理器直接使用
		c.Set("current_user_id", u.ID)

		// 提交事务（实际上只读，但保证原子性语义）
		if err := tx.Commit().Error; err != nil {
			// 提交失败概率极低，仍可返回 500
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// 通过检查，继续执行后续逻辑
		c.Next()
	}
}
