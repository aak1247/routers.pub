package v0

import (
	"github.com/gin-gonic/gin"
	"routers.pub/routers/middlewares"
)

func Init(r *gin.RouterGroup) {
	// 路由
	r.POST("/streams", middlewares.EnableTransaction, middlewares.WithRouterCtx(addStream))
	r.GET("/streams", middlewares.EnableTransaction, middlewares.WithRouterCtx(getStreams))
	r.GET("/streams/:streamId")
	r.PUT("/streams/:streamId", middlewares.EnableTransaction, middlewares.WithRouterCtx(updateStream))
	r.DELETE("/streams/:streamId", middlewares.EnableTransaction)
	r.PATCH("/streams/:streamId", middlewares.EnableTransaction)

	// web hook
	r.POST("/hooks/streams/:streamId", middlewares.EnableTransaction, middlewares.WithRouterCtx(callStream))
}
