# 介绍

## 功能

- Webhook通道，实现对请求的解析、请求的映射、请求的转发

## 核心概念

- `Stream`: 通道，定义了来源请求的方法和请求体格式、转发的地址和方法、转发的query、body和header等，并包含了来源请求中的字段
- `JSONSchema`: 用于定义转发请求体的格式，通过${variable_name} / #header.variable_name# / #query.variable_name#
  来引用来源请求中的变量（目前仅支持字符串）

## 使用场景

- 不同webhook的对接，直接定义转发的schema，通过schema实现对转发请求的定义，并获取到hook来源请求中的变量进行赋值
- 例如：报警系统与钉钉、企业微信、飞书等的对接，定义好转发的schema，通过schema实现对转发请求的定义，并获取到报警系统来源请求中的变量对转发请求中的占位符进行赋值
- 使用流程和工作模式：
    1. 创建一个Stream，定义好来源请求的方法、请求体格式、转发的地址和方法、转发的query、body和header等
    2. 在来源hook系统中配置webhook，将请求转发到Stream的地址
    3. stream解析请求，并根据定义生成请求，转发到stream中定义的地址

## API

### 创建Stream

#### URL

`POST /v0/streams`

#### 请求参数

| 名称                   | 位置   | 类型                        | 必选 | 说明                                                                        |
|----------------------|------|---------------------------|----|---------------------------------------------------------------------------|
| body                 | body | object                    | 否  | none                                                                      |
| » name               | body | string                    | 是  | none                                                                      |
| » description        | body | string                    | 是  | none                                                                      |
| » type               | body | string                    | 是  | none                                                                      |
| » requestContentType | body | string                    | 是  | 如application/json \| application/xml \| application/x-www-form-urlencoded |
| » hookBody           | body | [JSONSchema](#JSONSchema) | 是  | name为字段名  实际格式见schema定义和请求示例                                              |
| »» type              | body | string                    | 是  | none                                                                      |
| »» default           | body | string                    | 否  | none                                                                      |
| »» properties        | body | object                    | 否  | none                                                                      |
| »»» name             | body | [JSONSchema](#JSONSchema) | 是  | name为字段名 实际格式见schema定义和请求示例                                               |
| »»»» type            | body | string                    | 是  | none                                                                      |
| »»»» default         | body | string                    | 否  | none                                                                      |
| »»»» properties      | body | object                    | 否  | none                                                                      |
| » hookHeaders        | body | object                    | 是  | none                                                                      |
| »» key               | body | string                    | 是  | key / value 键值对                                                           |
| » hookParam          | body | [JSONSchema](#JSONSchema) | 是  | name为字段名                                                                  |
| »» type              | body | string                    | 是  | none                                                                      |
| »» default           | body | string                    | 否  | none                                                                      |
| »» properties        | body | object                    | 否  | none                                                                      |
| » hookUrl            | body | string                    | 是  | none                                                                      |
| » hookMethod         | body | string                    | 是  | none                                                                      |
| » mapping            | body | array<object>             | 是  | hookParam和requestParam的对象构成的数组，详见示例                                       |
| »» hookParam         | body | object                    | 是  | none                                                                      |
| »»» name             | body | string                    | 是  | none                                                                      |
| »»» type             | body | string                    | 是  | none                                                                      |
| »» requestParam      | body | object                    | 是  | none                                                                      |
| »»» name             | body | string                    | 是  | none                                                                      |
| »»» type             | body | string                    | 是  | none                                                                      |

#### 请求示例

创建一个stream，用于将sentry的报警转发到企业微信

```shell
curl --location --request POST 'http://routers.pub.aak1247.cn/v0/streams' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 118.195.177.165:8080' \
--header 'Connection: keep-alive' \
--data-raw '{
    "name": "Sentry template_card 测试",
    "description": "Sentry Webhook 测试",
    "type": "simple",
    "requestContentType": "application/json",
    "hookBody": {
        "type": "object",
        "properties": {
            "template_card": {
                "type": "object",
                "properties": {
                    "card_type": {
                        "type": "string",
                        "default": "text_notice"
                    },
                    "source": {
                        "type": "object",
                        "properties": {
                            "desc": {
                                "type": "string",
                                "default": "Sentry"
                            }
                        }
                    },
                    "main_title": {
                        "type": "object",
                        "properties": {
                            "title": {
                                "type": "string",
                                "default": "Sentry报警"
                            }
                        }
                    },
                    "horizontal_content_list": {
                        "type": "array",
                        "default": [
                            {
                                "keyname": "Project",
                                "value": "${project_name}"
                            },
                            {
                                "keyname": "预警级别",
                                "value": "${level}"
                            },
                            {
                                "keyname": "错误消息",
                                "value": "${message}"
                            },
                            {
                                "keyname": "原因",
                                "value": "${culprit}"
                            }
                        ]
                    },
                    "jump_list": {
                        "type": "array",
                        "default": [
                           {
                                "type": 1,
                                "url": "${url}",
                                "title": "Sentry报警详情"
                           }
                        ]
                    }
                }
            },
            "msgtype": {
                "type": "string",
                "default": "template_card"
            }
        }
    },
    "hookHeaders": {
    },
    "hookParam": {
        "type": "object",
        "properties": {
            "key": {
                "type": "string",
                "default": "xxxxxxxxxxxxxxxxxxxxxxxxx"
            }
        }
    },
    "hookContentType": "application/json",
    "hookUrl": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send",
    "hookMethod": "POST",
    "mapping": [
        {
            "hookParam": {
                "name": "template_card.source.desc",
                "type": "body"
            },
            "requestParam": {
                "name": "content",
                "type": "body"
            }
        }
    ]
}'
```

#### 响应示例

响应中的`id`即为stream的id，用于后续的请求中

```json
{
  "id": "2e6ed59b-4ef3-44b7-9839-65466a93a0b5",
  "createdAt": "2023-08-24T07:45:22.07131292Z",
  "updatedAt": "2023-08-24T07:45:22.07131292Z",
  "name": "Sentry template_card 测试",
  "description": "Sentry Webhook 测试",
  "type": "simple",
  "requestContentType": "application/json",
  "hookBody": {
    "type": "object",
    "default": null,
    "properties": {
      "msgtype": {
        "type": "string",
        "default": "template_card",
        "properties": null
      },
      "template_card": {
        "type": "object",
        "default": null,
        "properties": {
          "card_type": {
            "type": "string",
            "default": "text_notice",
            "properties": null
          },
          "horizontal_content_list": {
            "type": "array",
            "default": [
              {
                "keyname": "Project",
                "value": "${project_name}"
              },
              {
                "keyname": "预警级别",
                "value": "${level}"
              },
              {
                "keyname": "错误消息",
                "value": "${message}"
              },
              {
                "keyname": "原因",
                "value": "${culprit}"
              }
            ],
            "properties": null
          },
          "jump_list": {
            "type": "array",
            "default": [
              {
                "title": "Sentry报警详情",
                "type": 1,
                "url": "${url}"
              }
            ],
            "properties": null
          },
          "main_title": {
            "type": "object",
            "default": null,
            "properties": {
              "title": {
                "type": "string",
                "default": "Sentry报警",
                "properties": null
              }
            }
          },
          "source": {
            "type": "object",
            "default": null,
            "properties": {
              "desc": {
                "type": "string",
                "default": "Sentry",
                "properties": null
              }
            }
          }
        }
      }
    }
  },
  "hookHeaders": {},
  "hookParam": {
    "type": "object",
    "default": null,
    "properties": {
      "key": {
        "type": "string",
        "default": "xxxxxxxxxxxxxxxxxxxxx",
        "properties": null
      }
    }
  },
  "hookUrl": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send",
  "hookMethod": "POST",
  "hookContentType": "application/json",
  "mapping": [
    {
      "hookParam": {
        "name": "text.content",
        "type": "body"
      },
      "requestParam": {
        "name": "content",
        "type": "body"
      }
    }
  ]
}
```

### 调用stream

#### 请求示例

> url: `http://routers.pub.aak1247.cn/v0/hooks/streams/:streamId`

将streamId替换为上一步中获取到的id

```bash
curl --location --request POST 'http://routers.pub.aak1247.cn/v0/hooks/streams/2e6ed59b-4ef3-44b7-9839-65466a93a0b5' \
--header 'User-Agent: Apifox/1.0.0 (https://apifox.com)' \
--header 'Content-Type: application/json' \
--header 'Accept: */*' \
--header 'Host: 118.195.177.165:8080' \
--header 'Connection: keep-alive' \
--data-raw '{
    "project_name": "test project",
    "message": "测试预警",
    "level": "error",
    "culprit": "测试",
    "content": "内容",
    "url": "https://sentry.io/organizations/xxxx/issues/xxxxx/?project=xxxxx"
}'
```

### 更新stream

> PUT /v0/streams/:streamId

请求体与创建stream时的请求体相同

## 结构定义

### JSONSchema

#### 属性

| 名称         | 类型                                             | 必选    | 约束   | 中文名 | 说明                                 |
|------------|------------------------------------------------|-------|------|-----|------------------------------------|
| type       | string                                         | true  | none |     | string/object/number/boolean/array |
| default    | string \| object \| boolean \| array \| number | false | none |     | none                               |
| properties | object                                         | false | none |     | none                               |
| » name     | [JSONSchema](#JSONSchema)                      | true  | none |     | name为字段名                           |

#### 示例

```json
{
  "type": "object",
  "properties": {
    "property1": {
      "type": "object",
      "default": {},
      "properties": {
        "name": {
          "type": "string",
          "default": "test"
        }
      }
    }
  }
}
```
