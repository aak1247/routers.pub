package middlewares

import (
	"github.com/gin-gonic/gin"
	"routers.pub/infra"
)

func Secure(c *gin.Context) {
	// 检查请求是否有效（access token）
	// TODO
	// 初始化请求id
	requestId := c.Request.Header.Get("X-Request-Id")
	if requestId == "" {
		requestId = infra.NewIdStr()
	}
	c.Set("requestId", requestId)
}
