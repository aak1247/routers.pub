package domains

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"routers.pub/framework"
	"routers.pub/infra"
	"routers.pub/utils"
	"strings"
)

type (
	StreamType  = string
	HttpMethod  = string
	ContentType = string

	ParamType = string

	Param struct {
		// Name
		Name string `json:"name"`
		// Type
		Type ParamType `json:"type"`
	}

	// Stream
	Stream struct {
		Entity
		utils.CanHasError `gorm:"-" json:"-"`
		// Name
		Name string `gorm:"column:name;not null" json:"name" search:"key:name"`
		// Description
		Description string `gorm:"column:description" json:"description"`
		// Type
		Type StreamType `gorm:"column:type;not null" json:"type" search:"key:type"`
		// 请求体格式
		RequestContentType ContentType `gorm:"column:request_content_type;not null" json:"requestContentType"`
		// hook体
		HookBody *utils.JSONSchema `gorm:"column:hook_body;not null" json:"hookBody"`
		// hook头
		HookHeaders *utils.StringMap `gorm:"column:hook_headers;not null" json:"hookHeaders"`
		// hook参数
		HookParam *utils.JSONSchema `gorm:"column:hook_param;not null" json:"hookParam"`
		// hook url
		HookUrl string `gorm:"column:hook_url;not null" json:"hookUrl" search:"key:hookUrl"`
		// hook method
		HookMethod HttpMethod `gorm:"column:hook_method;not null" json:"hookMethod"`
		// hook type
		HookContentType ContentType `gorm:"column:hook_content_type" json:"hookContentType"`
		// 映射关系 key: request param， value: hook param
		Mapping *ParamMap `gorm:"column:mapping;not null" json:"mapping"`

		// 下为状态字段
		// 实际请求url
		ActualRequestUrl string `gorm:"-" json:"-"`
		// 实际请求头
		ActualRequestHeaders *map[string]string `gorm:"-" json:"-"`
		// 实际请求参数
		ActualRequestParams *map[string]string `gorm:"-" json:"-"`
		// 实际请求体
		ActualRequestBody *map[string]interface{} `gorm:"-" json:"-"`
		// 实际hook头
		ActualHookHeaders *map[string]string      `gorm:"-" json:"-"`
		ActualHookParams  *map[string]interface{} `gorm:"-" json:"-"`
		// 实际hook体
		ActualHookBody *map[string]interface{} `gorm:"-" json:"-"`
		// 缓存所有变量（带路径），一级
		AllVariables *map[string]interface{} `gorm:"-" json:"-"`
	}
)
type Streams struct {
	utils.CanHasError `gorm:"-" json:"-"`
	Streams           []*Stream `json:"streams"`
}

const (
	StreamTypeSimple StreamType = "simple"
	StreamTypeScript StreamType = "script"

	HttpMethodGet     HttpMethod = "GET"
	HttpMethodPost    HttpMethod = "POST"
	HttpMethodPut     HttpMethod = "PUT"
	HttpMethodDelete  HttpMethod = "DELETE"
	HttpMethodPatch   HttpMethod = "PATCH"
	HttpMethodOptions HttpMethod = "OPTIONS"

	ContentTypeJson      ContentType = "application/json"
	ContentTypeForm      ContentType = "application/x-www-form-urlencoded"
	ContentTypeText      ContentType = "text/plain"
	ContentTypeHtml      ContentType = "text/html"
	ContentTypeXml       ContentType = "text/xml"
	ContentTypeMultipart ContentType = "multipart/form-data"

	ParamTypeHeader ParamType = "header"
	ParamTypeQuery  ParamType = "query"
	ParamTypePath   ParamType = "path"
	ParamTypeBody   ParamType = "body"
)

var (
	HEADER_MATCHER = utils.NewMatcher("#header\\.([^#]+)#")
	QUERY_MATCHER  = utils.NewMatcher("#query\\.([^#]+)#")
	BODY_MATCHER   = utils.NewMatcher("\\${([^}]+)}")
)

// TableName sets the insert table name for this struct type
func (h *Stream) TableName() string {
	return "streams"
}

func (h *Stream) SetId(id string) *Stream {
	h.ID = id
	return h
}

func (h *Stream) UpdateByParam(param interface{}) *Stream {
	t := utils.UpdateEntityByDef(*h, param)
	return &t
}

func (h *Stream) ParseRequest(req *http.Request) *Stream {
	if h.HasError() {
		return h
	}
	// 解析请求头
	headers := make(map[string]string)
	for key, value := range req.Header {
		headers[key] = value[0]
	}
	h.ActualRequestHeaders = &headers
	infra.Log.Debug("request headers: ", headers)

	// 解析请求参数
	params := make(map[string]string)
	for key, value := range req.URL.Query() {
		params[key] = value[0]
	}
	h.ActualRequestParams = &params
	infra.Log.Debug("request params: ", params)

	// 解析请求体
	body := make(map[string]interface{})
	switch h.RequestContentType {
	case ContentTypeJson:
		_ = json.NewDecoder(req.Body).Decode(&body)
	case ContentTypeForm:
		_ = req.ParseForm()
		for key, value := range req.Form {
			body[key] = value[0]
		}
	case ContentTypeXml:
		_ = xml.NewDecoder(req.Body).Decode(&body)
	default:
		err := framework.NewError("不支持的请求体格式")
		h.AddError(err)
	}
	h.ActualRequestBody = &body
	// 日志
	infra.Log.Debug("request body: ", body)

	// 缓存所有变量到allVariables
	allVariables := make(map[string]interface{})
	for key, value := range headers {
		allVariables["#header."+key+"#"] = value
	}
	for key, value := range params {
		allVariables["#query."+key+"#"] = value
	}
	flattenBody := utils.FlattenMap(body, ".", "")
	for key, value := range flattenBody {
		allVariables["${"+key+"}"] = value
	}
	h.AllVariables = &allVariables
	infra.Log.Debug("request variables: ", allVariables)

	return h
}

func (h *Stream) Generate() *Stream {
	if h.HasError() {
		return h
	}

	// 根据hook headers、params、body 定义，生成实际的请求头、参数、体
	// header
	headers := make(map[string]string)
	if h.HookHeaders != nil {
		for k, v := range *h.HookHeaders {
			headers[k] = v
		}
	}
	h.ActualHookHeaders = &headers
	// param
	if params, ok := h.HookParam.GetDefaultValue().(map[string]interface{}); ok {
		h.ActualHookParams = &params
	}

	// body
	if body, ok := h.HookBody.GetDefaultValue().(map[string]interface{}); ok {
		h.ActualHookBody = &body
	}
	return h
}

func (h *Stream) Transform() *Stream {
	if h.HasError() {
		return h
	}

	// 根据mapping，转换请求参数
	if h.Mapping != nil {
		reqMap := make(map[ParamType]interface{})
		reqMap[ParamTypeBody] = h.ActualRequestBody
		reqMap[ParamTypeHeader] = h.ActualRequestHeaders
		reqMap[ParamTypeQuery] = h.ActualRequestParams
		reqMap[ParamTypePath] = make(map[string]string) // TODO：目前不支持path参数
		hookMap := make(map[ParamType]interface{})
		hookMap[ParamTypeBody] = h.ActualHookBody
		hookMap[ParamTypeHeader] = h.ActualHookHeaders
		hookMap[ParamTypeQuery] = h.ActualHookParams
		hookMap[ParamTypePath] = make(map[string]string) // TODO：目前不支持path参数
		for _, mapping := range *h.Mapping {
			hk, req := mapping.HookParam, mapping.RequestParam
			hkMp := hookMap[hk.Type]
			rqMp := reqMap[req.Type]
			if hkMp == nil || rqMp == nil {
				continue
			}
			var targetValue interface{}
			if rqMapInterface, ok := rqMp.(map[string]interface{}); ok {
				targetValue = utils.DepGetFromMap(rqMapInterface, req.Name, ".")
			} else if rqMapString, ok := rqMp.(map[string]string); ok {
				targetValue = rqMapString[req.Name]
			}
			if targetValue == nil {
				continue
			}
			if hkMapInterface, ok := hkMp.(map[string]interface{}); ok {
				utils.DepSetToMap(hkMapInterface, hk.Name, ".", targetValue)
			} else if hkMapString, ok := hkMp.(map[string]string); ok {
				targetValueString := fmt.Sprintf("%v", targetValue)
				hkMapString[hk.Name] = targetValueString
			}
		}
	}

	// 处理hook各字段中的占位符
	if h.ActualHookHeaders != nil {
		utils.TraverseStringMap(*h.ActualHookHeaders, func(key string, value string) {
			(*h.ActualHookHeaders)[key] = h.formatValue(value)
		})
	}

	if h.ActualHookParams != nil {
		utils.TraverseMapString(*h.ActualHookParams, func(key string, value string) string {
			return h.formatValue(value)
		})
	}

	if h.ActualHookBody != nil {
		utils.TraverseMapString(*h.ActualHookBody, func(key string, value string) string {
			return h.formatValue(value)
		})
	}

	return h
}

func (h *Stream) formatValue(value string) string {
	// 匹配占位符
	headerMatched := HEADER_MATCHER.FindAllStringSubmatch(value, -1)
	queryMatched := QUERY_MATCHER.FindAllStringSubmatch(value, -1)
	bodyMatched := BODY_MATCHER.FindAllStringSubmatch(value, -1)
	allMatched := append(append(headerMatched, queryMatched...), bodyMatched...)
	for _, matched := range allMatched {
		matchText := matched[0]
		// 获取变量名
		//variableName := matched[1] // 这里得测一下
		// 获取变量值
		variableValue := (*h.AllVariables)[matchText]
		// 替换占位符
		value = strings.ReplaceAll(value, matchText, fmt.Sprintf("%v", variableValue))
	}
	return value
}

func (h *Stream) DoRequest() (*http.Response, error) {

	if h.HasError() {
		return nil, h.GetError()
	}

	// 根据hook方法、url、headers、params、body，发起请求
	client := http.Client{}
	// 对param进行url编码
	url := h.HookUrl
	if h.ActualHookParams != nil {
		url += "?" + utils.UrlEncode(*h.ActualHookParams)
	}
	var body []byte
	switch h.HookContentType {
	case ContentTypeJson:
		body, _ = json.Marshal(h.ActualHookBody)
		h.HookHeaders.Set("Content-Type", "application/json")
	case ContentTypeForm:
		body = []byte(utils.UrlEncode(*h.ActualHookBody))
		h.HookHeaders.Set("Content-Type", "application/x-www-form-urlencoded")
	case ContentTypeXml:
		body, _ = xml.Marshal(h.ActualHookBody)
		h.HookHeaders.Set("Content-Type", "application/xml")
	}
	req, err := http.NewRequest(h.HookMethod, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	// header
	if h.ActualHookHeaders != nil {
		for key, value := range *h.ActualHookHeaders {
			req.Header.Add(key, value)
		}
	}
	return client.Do(req)
}

func (h *Stream) Find(ctx *framework.RouterCtx) *Stream {
	stream := framework.Cached(ctx, h.GetCacheKey(), func() interface{} {
		stream := &Stream{}
		if err := ctx.GetDb().Model(h).Where("id = ?", h.ID).First(stream).Error; err != nil {
			stream.AddErrors(h.GetErrors())
			stream.AddError(err)
		}
		return stream
	}).(*Stream)
	return stream
}

func (h *Stream) FindAll(ctx *framework.RouterCtx) *Streams {
	streams := framework.Cached(ctx, framework.GetSearchKey(h, "streams:"), func() interface{} {
		streams := &Streams{}
		streamList := make([]*Stream, 0)
		if err := ctx.GetDb().Model(h).Find(&streamList).Error; err != nil {
			streams.AddErrors(h.GetErrors())
			streams.AddError(err)
		}
		streams.Streams = streamList
		return streams
	}).(*Streams)
	return streams
}

func (h *Stream) Save(ctx *framework.RouterCtx) *Stream {
	ctx.SetToLocalCache(h.GetCacheKey(), h)
	if err := ctx.GetDb().Save(h).Error; err != nil {
		h.AddError(err)
	}
	return h
}

func (h *Stream) Delete(ctx *framework.RouterCtx) *Stream {
	ctx.DelFromLocalCache(h.GetCacheKey())
	if err := ctx.GetDb().Delete(h).Error; err != nil {
		h.AddError(err)
	}
	return h
}

// utils
type ParamMap []struct {
	HookParam    Param `json:"hookParam"`
	RequestParam Param `json:"requestParam"`
}

func (pm *ParamMap) Value() (driver.Value, error) {
	return json.Marshal(pm)
}

// Scan 实现方法
func (pm *ParamMap) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &pm)
}

func init() {
	registerAutoMigrate(&Stream{})
}
