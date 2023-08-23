package framework

import (
	"github.com/gin-gonic/gin"
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

func NewRouterCtx() *RouterCtx {
	return &RouterCtx{
		db:    dbs.GetDb(),
		cache: make(map[string]interface{}, perReqMaxCacheSize),
	}
}

func (ctx *RouterCtx) GetDb() *gorm.DB {
	return ctx.db
}

func (ctx *RouterCtx) EnableTransaction() {
	ctx.transactionEnabled = true
	ctx.db = dbs.StartTx()
}

func (ctx *RouterCtx) SetGinCtx(c *gin.Context) {
	ctx.ctx = c
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

func (ctx *RouterCtx) SetToLocalCache(key string, value interface{}) {
	cacheKey := key
	//if !globalKey {
	//	getCacheKey(key, 2)
	//}
	reqCache := ctx.cache
	reqCache[cacheKey] = value
}

func (ctx *RouterCtx) DelFromLocalCache(key string) {
	cacheKey := key
	//if !globalKey {
	//	getCacheKey(key, 2)
	//}
	reqCache := ctx.cache
	delete(reqCache, cacheKey)
}
