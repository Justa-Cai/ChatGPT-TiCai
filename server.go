package main

import (
	"github.com/gin-gonic/gin"
)

func start_http_server() {
	r := gin.Default()

	// 设置静态文件路由和目录
	r.Static("/tmp", "/tmp/")

	r.Run(":61234")
}
