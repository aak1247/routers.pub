package v0

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"routers.pub/domains"
	"routers.pub/framework"
	"routers.pub/utils"
)

type (
	AddStreamReq struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Type        string `json:"type" binding:"required"`
		// 请求体格式
		RequestContentType string            `json:"requestContentType" binding:"required"`
		HookContentType    string            `json:"hookContentType" binding:"required"`
		HookBody           *utils.JSONSchema `json:"hookBody"`
		HookHeaders        *utils.StringMap  `json:"hookHeaders"`
		HookParam          *utils.JSONSchema `json:"hookParam"`
		HookUrl            string            `json:"hookUrl" binding:"required"`
		HookMethod         string            `json:"hookMethod" binding:"required"`
		Mapping            *domains.ParamMap `json:"mapping"`
	}
)

func callStream(c *gin.Context, routerCtx *framework.RouterCtx) {
	streamId := c.Param("streamId")
	var stream = (&domains.Stream{}).SetId(streamId)
	if stream = stream.Find(routerCtx); stream.HasError() {
		routerCtx.AddError(stream.GetError())
		return
	}

	resp, err := stream.ParseRequest(c.Request). // 根据stream parse请求
							Generate().  // 根据stream generate请求
							Transform(). // 根据stream transform请求
							DoRequest()  // 发送hook请求
	if err != nil {
		routerCtx.AddError(err)
	} else {
		var bodyText []byte
		bodyText, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			routerCtx.AddError(err)
		}
		_, err = c.Writer.Write(bodyText)
		if err != nil {
			routerCtx.AddError(err)
		} else {
			c.Writer.WriteHeader(resp.StatusCode)
		}
	}
	return
}

// 创建stream
func AddStream(c *gin.Context, routerCtx *framework.RouterCtx) {
	var req AddStreamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		routerCtx.AddError(err)
		return
	}
	stream := (&domains.Stream{}).UpdateByParam(req)
	if stream.HasError() {
		routerCtx.AddError(stream.GetError())
		return
	}
	if stream = stream.Save(routerCtx); stream.HasError() {
		routerCtx.AddError(stream.GetError())
		return
	}
	c.JSON(http.StatusOK, stream)
}

func UpdateStream(c *gin.Context, routerCtx *framework.RouterCtx) {
	var req AddStreamReq
	if err := c.ShouldBindJSON(&req); err != nil {
		routerCtx.AddError(err)
		return
	}
	id := c.Param("streamId")
	if id == "" {
		routerCtx.AddError(framework.NewError("id is empty"))
		return
	}
	stream := (&domains.Stream{}).SetId(id).
		Find(routerCtx).
		UpdateByParam(req).
		Save(routerCtx)
	if stream.HasError() {
		routerCtx.AddError(stream.GetError())
		return
	}
	c.JSON(http.StatusOK, stream)
}
