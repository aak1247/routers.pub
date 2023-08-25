package framework

import (
	"errors"
	"routers.pub/utils"
)

func NewError(msg string) error {
	return errors.New(msg)
}

func AddErrorTo[T any](err error, maybe T) T {
	// 检查是否是ICanAddError接口
	var i interface{}
	i = maybe
	if maybe, ok := i.(utils.ICanHasError); ok {
		maybe.AddError(err)
		return maybe.(T)
	}
	return maybe
}
