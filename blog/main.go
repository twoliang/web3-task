package main

import (
	"web-task/blog/utility"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	utility.InitDB()

	// 创建Gin实例
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	r.Run(":8080") // 默认监听8080端口
}
