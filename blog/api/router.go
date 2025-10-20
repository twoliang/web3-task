// api/router.go
package api

import (
	"web-task/blog/internal/controller"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化所有路由
func SetupRouter(
	r *gin.Engine,
	userCtl *controller.UserController,
	postCtl *controller.PostHandler,
	commentCtl *controller.CommentHandler,
) {
	// 基础API前缀
	api := r.Group("/api/v1")

	// 注册各模块路由
	SetupUserRouter(api, userCtl)       // 用户路由
	SetupPostRouter(api, postCtl)       // 文章路由
	SetupCommentRouter(api, commentCtl) // 评论路由
}
