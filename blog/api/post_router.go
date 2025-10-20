// api/post_router.go
package api

import (
	"web-task/blog/internal/controller"

	"github.com/gin-gonic/gin"
)

// SetupPostRouter 注册文章相关路由
func SetupPostRouter(router *gin.RouterGroup, pc *controller.PostHandler) {
	// 直接调用处理器自带的路由注册方法
	pc.RegisterRoutes(router)
}
