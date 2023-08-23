package middlewares

import (
	"github.com/gin-gonic/gin"
	"time"
)

func Profile(c *gin.Context) {
	// 记录开始时间
	c.Set("x-start-time", time.Now().UnixMilli())
}
