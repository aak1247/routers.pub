package framework

import (
	"reflect"
	"routers.pub/utils"
)

type cachedValues struct {
	Values []interface{}
}

// 请求级cache, 用于存储请求级别的数据
var (
	perReqMaxCacheSize = 1000
)

func getCacheKey(key string, callDepth int) string {
	prefix := ""
	//pc, file, _, ok := runtime.Caller(callDepth)
	//if ok {
	//	pcName := runtime.FuncForPC(pc).Name() //获取函数名
	//	prefix = fmt.Sprintf("%s:%s:", file, pcName)
	//}
	return prefix + key
}

// getter 是返回cache value的函数
// returns cache value, 目前只能接受单个返回值，多返回值会有类型问题（返回类型为[]interface{}）
func Cached(ctx *RouterCtx, key string, fn interface{}, args ...interface{}) interface{} {
	var getter = func() interface{} {
		// 利用反射调用函数f并将args作为参数
		fnValue := reflect.ValueOf(fn)
		argsValue := utils.Map(args, func(i int, arg interface{}) reflect.Value {
			return reflect.ValueOf(arg)
		})
		values := utils.Map(fnValue.Call(argsValue), func(i int, value reflect.Value) interface{} {
			return value.Interface()
		})
		return &cachedValues{Values: values}
	}
	res := ctx.cacheable(key, getter)
	if resValue, ok := res.(*cachedValues); ok {
		// 拆包
		if len(resValue.Values) == 1 {
			return resValue.Values[0]
		} else {
			return resValue.Values
		}
	} else {
		return res
	}
}
