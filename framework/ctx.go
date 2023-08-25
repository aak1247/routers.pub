package framework

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"routers.pub/dbs"
	"routers.pub/utils"
)

type (
	RouterCtx struct {
		utils.CanHasError
		ctx                *gin.Context
		RequestId          string
		db                 *gorm.DB
		transactionEnabled bool
		cache              map[string]interface{}
	}
)

var (
	myValidater = validator.New()
)

func NewRouterCtx() *RouterCtx {
	return &RouterCtx{
		db:    dbs.GetDb(),
		cache: make(map[string]interface{}, perReqMaxCacheSize),
	}
}

func (ctx *RouterCtx) GetDb() *gorm.DB {
	return ctx.db
}

func (ctx *RouterCtx) EnableTransaction() *RouterCtx {
	ctx.transactionEnabled = true
	ctx.db = dbs.StartTx()
	return ctx
}

func (ctx *RouterCtx) SetGinCtx(c *gin.Context) *RouterCtx {
	ctx.ctx = c
	return ctx
}

func (ctx *RouterCtx) Release() error {
	// 清空缓存
	ctx.cache = nil
	if !ctx.transactionEnabled {
		return nil
	}
	if ctx.HasError() {
		return ctx.db.Rollback().Error
	}
	return ctx.db.Commit().Error
}

func (ctx *RouterCtx) GetFromLocalCache(key string) interface{} {
	cacheKey := key
	//if !globalKey {
	//	getCacheKey(key, 2)
	//}
	reqCache := ctx.cache
	return reqCache[cacheKey]
}

func (ctx *RouterCtx) cacheable(key string, getter func() interface{}) interface{} {
	cacheKey := key
	//if !globalKey {
	//	getCacheKey(key, 2)
	//}
	reqCache := ctx.cache
	if v, ok := reqCache[cacheKey]; ok {
		return v
	}
	v := getter()
	reqCache[cacheKey] = v
	return v
}

func (ctx *RouterCtx) SetToLocalCache(key string, value interface{}) *RouterCtx {
	cacheKey := key
	//if !globalKey {
	//	getCacheKey(key, 2)
	//}
	reqCache := ctx.cache
	reqCache[cacheKey] = value
	return ctx
}

func (ctx *RouterCtx) DelFromLocalCache(key string) *RouterCtx {
	cacheKey := key
	//if !globalKey {
	//	getCacheKey(key, 2)
	//}
	reqCache := ctx.cache
	delete(reqCache, cacheKey)
	return ctx
}

func (ctx *RouterCtx) BindQuery(resp interface{}) *RouterCtx {
	err := bindQuery(ctx.ctx, resp)
	if err != nil {
		ctx.AddError(err)
		AddErrorTo(err, resp)
	}
	// validate
	err = myValidater.Struct(resp)
	if err != nil {
		ctx.AddError(err)
		AddErrorTo(err, resp)
	}
	return ctx
}

func (ctx *RouterCtx) BindJSON(resp interface{}) *RouterCtx {
	err := ctx.ctx.ShouldBindJSON(resp)
	if err != nil {
		ctx.AddError(err)
		AddErrorTo(err, resp)
	}
	// validate
	err = myValidater.Struct(resp)
	if err != nil {
		ctx.AddError(err)
		AddErrorTo(err, resp)
	}
	return ctx
}

func (ctx *RouterCtx) BindParam(name string, param *string, required bool) *RouterCtx {
	p := ctx.ctx.Param(name)
	if p == "" && required {
		ctx.AddError(NewError("param " + name + " is required"))
		return ctx
	}
	*param = p
	return ctx
}

func (ctx *RouterCtx) Response(resp interface{}) *RouterCtx {
	if typedResp, ok := resp.(utils.ICanHasError); ok {
		if typedResp.HasError() {
			ctx.AddErrors(typedResp.GetErrors())
		}
	}
	if ctx.HasError() {
		msg := ctx.AsErrorMessage()
		ctx.ctx.JSON(msg.StatusCode, msg)
		return ctx
	}
	switch resp.(type) {
	case string:
		ctx.ctx.String(200, resp.(string))
	case []byte:
		ctx.ctx.Data(200, "application/octet-stream", resp.([]byte))
	}
	ctx.ctx.JSON(200, resp)
	return ctx
}
