package main

import (
	"routers.pub/dbs"
	"routers.pub/domains"
	"routers.pub/env"
	"routers.pub/infra"
	"routers.pub/routers"
)

func main() {
	env.Conf = env.GetEnv()
	infra.InitLogger()
	infra.InitId()
	dbs.InitDatabase()
	if err := domains.InitDbTables(); err != nil {
		panic(err)
	}
	if err := dbs.InitRedis(); err != nil {
		panic(err)
	}
	routers.Init()
}
