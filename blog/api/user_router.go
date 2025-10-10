package api

import (
	"web-task/blog/internal/controller"

	"github.com/gin-gonic/gin"
)

// 定义用户相关路由
func SetupUserRouter(router *gin.RouterGroup, uc *controller.UserController) {
	userRouter := router.Group("/users")
	{
		userRouter.POST("/register", uc.Register) // 用户注册
		userRouter.POST("/login", uc.Login)       // 用户登录
	}
}
