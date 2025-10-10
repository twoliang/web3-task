// router/router.go
package api

import (
	"web-task/blog/internal/controller"

	"github.com/gin-gonic/gin"

	// 引入你的中间件（如JWT鉴权）
	"web-task/blog/middleware"
)

// InitRouter 初始化所有路由（总入口）
// 参数：接收所有模块的控制器实例（通过依赖注入传递）
func InitRouter(
	r *gin.Engine,
	userCtrl *controller.UserController,
) {
	// 1. 注册全局中间件（对所有请求生效）
	r.Use(gin.Logger())      // 日志中间件
	r.Use(gin.Recovery())    // 异常恢复中间件
	r.Use(middleware.CORS()) // 跨域中间件（自定义）

	// 2. 创建API根路由组（所有接口统一前缀 /api/v1）
	apiV1 := r.Group("/api/v1")

	// 3. 注册各模块的路由（将根路由组传递给各模块）
	SetupUserRouter(apiV1, userCtrl) // 注册用户模块路由
}
