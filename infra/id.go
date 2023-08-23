package infra

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

var (
	startTime       = "2023-08-22"
	machineID int64 = 1
	node      *snowflake.Node
)

func InitId() {
	var st time.Time
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		fmt.Println(err)
		return
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func NewId() snowflake.ID {
	return node.Generate()
}

func NewIdStr() string {
	return node.Generate().String()
}
