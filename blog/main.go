package main

import (
	"web-task/blog/api"

	"web-task/blog/internal/controller"
	"web-task/blog/internal/logic"
	"web-task/blog/utility"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	db := utility.DB

	// 初始化控制器
	userService := logic.NewUserService(db)
	userCtl := controller.NewUserController(userService)

	postService := logic.NewPostService(db)
	postCtl := controller.NewPostHandler(postService)

	commentService := logic.NewCommentService(db)
	commentCtl := controller.NewCommentHandler(commentService)

	// 初始化gin引擎
	r := gin.Default()

	// 注册所有路由
	api.SetupRouter(r, userCtl, postCtl, commentCtl)

	// 启动服务
	r.Run(":8080")
}
