package routers

import (
	"github.com/gin-gonic/gin"
	"routers.pub/env"
	"routers.pub/infra"
	"routers.pub/routers/middlewares"
	v0 "routers.pub/routers/v0"
)

func Init() {
	if env.Conf.Server.Mode == "prod" || env.Conf.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(
		middlewares.Profile,
		middlewares.ErrorRecovery,
		middlewares.Cors(),
		middlewares.Secure,
		middlewares.InitCtx,
	)
	// v0
	groupV0 := r.Group("/v0")
	v0.Init(groupV0)

	err := r.Run(env.GetEnv().GetAddr()) // listen and serve on 0.0.0.0:8080}
	if err != nil {
		infra.Log.Fatal(err)
	}
}
