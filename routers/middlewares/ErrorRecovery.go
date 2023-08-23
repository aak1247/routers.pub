package middlewares

import (
	"github.com/gin-gonic/gin"
	"routers.pub/infra"
	"routers.pub/utils"
)

func ErrorRecovery(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			infra.Log.Errorf("[ERROR-ALERT] api panic: %v\n Stack:%s", r, utils.CallStack(20, 1))
			//封装通用json返回
			//c.JSON(http.StatusOK, Result.Fail(errorToString(r)))
			// TODO
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			c.Abort()
		}
	}()
}
