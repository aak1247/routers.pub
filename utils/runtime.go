package utils

import (
	"fmt"
	"runtime"
)

// CallStack
//
//	@Description: 读取调用堆栈字符串
//	@param line 读取最大行数
//	@param skip 跳过行数
//	@return string
func CallStack(line int, skip int) string {
	skip += 1
	str := "\n"
	for i := 0; i < line; i++ {
		pc, file, lineNumber, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		skip++
		pcName := runtime.FuncForPC(pc).Name() //获取函数名
		str = str + fmt.Sprintf("\t %s:%d[%s]\n", file, lineNumber, pcName)
	}
	return str
}
