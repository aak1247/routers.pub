package middlewares

import (
	"github.com/gin-gonic/gin"
	"routers.pub/framework"
	"routers.pub/infra"
)

func InitCtx(c *gin.Context) {
	routerCtx := framework.NewRouterCtx()
	routerCtx.SetGinCtx(c)
	// 获取请求id
	requestId, _ := c.Get("requestId")
	routerCtx.RequestId = requestId.(string)
	c.Set("ctx", routerCtx)
	// 完成后续接口调用
	c.Next()
	// 释放资源
	err := routerCtx.Release()
	if err != nil {
		infra.AlertError(err)
		routerCtx.AddError(err)
	}
	if routerCtx.HasError() {
		c.AbortWithStatusJSON(400, gin.H{
			"error": routerCtx.GetError().Error(),
		})
	}
}

func EnableTransaction(c *gin.Context) {
	routerCtx, ok := c.Get("ctx")
	if !ok {
		infra.AlertMessage("ctx not found")
		return
	}
	routerCtx.(*framework.RouterCtx).EnableTransaction()
}

func WithRouterCtx(f func(c *gin.Context, ctx *framework.RouterCtx)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, ok := c.Get("ctx")
		if !ok {
			infra.AlertMessage("ctx not found")
			return
		}
		f(c, ctx.(*framework.RouterCtx))
	}
}
