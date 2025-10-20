package api

import (
	"web-task/blog/internal/controller"

	"github.com/gin-gonic/gin"
)

// SetupCommentRouter 注册评论相关路由
func SetupCommentRouter(router *gin.RouterGroup, cc *controller.CommentHandler) {
	// 直接调用处理器自带的路由注册方法
	cc.RegisterRoutes(router)
}
