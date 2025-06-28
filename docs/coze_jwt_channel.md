# Coze JWT 渠道使用文档

## 概述

Coze JWT 渠道是一个新增的渠道类型，支持通过 OAuth JWT 授权（开发者）方式请求 Coze 工作流和智能体。该渠道基于官方 `github.com/coze-dev/coze-go` SDK 实现，使用 RSA 私钥签名的 JWT token 来获取 OAuth access token，提供更安全、更可靠的认证方式。

## 功能特性

### 🎯 核心功能
- **官方SDK支持**：基于 `github.com/coze-dev/coze-go` 官方SDK实现
- **OAuth JWT 认证**：使用 RSA 私钥签名的 JWT token 进行认证
- **智能体支持**：支持调用 Coze 智能体进行对话
- **同步工作流**：支持同步工作流执行，自动检测执行模式
- **异步工作流**：支持长时间运行的异步工作流，自动轮询状态
- **Token 缓存**：自动缓存和刷新 access token
- **流式响应**：支持流式和非流式响应
- **状态轮询**：异步工作流自动轮询执行状态直到完成

### 📊 支持的请求类型
| 类型 | 模型名称格式 | 说明 | API 端点 | 执行模式 |
|------|-------------|------|----------|----------|
| 智能体 | 普通模型名 | 用于对话聊天 | `/v3/chat` | 异步轮询 |
| 同步工作流 | `workflow:{workflow_id}` | 用于工作流执行 | `/v1/workflow/run` | 自动检测 |
| 异步工作流 | `workflow-async:{workflow_id}` | 用于长时间工作流 | `/v1/workflow/run` | 状态轮询 |

## 配置说明

### 1. 渠道基本信息
- **渠道类型**：选择 "Coze JWT" (类型值: 51)
- **渠道名称**：自定义名称
- **Base URL**：`https://api.coze.cn`
- **模型列表**：支持的模型名称，用逗号分隔

### 2. JWT 配置 (渠道额外设置字段)
在渠道的"渠道额外设置"字段中配置 JSON 格式的认证信息：

```json
{
  "client_id": "你的客户端ID",
  "public_key_id": "你的公钥ID", 
  "private_key": "-----BEGIN PRIVATE KEY-----\n你的RSA私钥\n-----END PRIVATE KEY-----",
  "space_id": "你的空间ID",
  "default_bot_id": "默认智能体ID"
}
```

#### 配置字段说明
- `client_id`: Coze 应用的客户端 ID
- `public_key_id`: 公钥 ID，用于 JWT header 中的 kid 字段
- `private_key`: RSA 私钥，用于签名 JWT token
- `space_id`: Coze 空间 ID（可选）
- `default_bot_id`: 默认智能体 ID，当请求中未指定时使用

### 3. 测试参数示例
```json
{
  "client_id": "1123962302922",
  "public_key_id": "fxkr4uRho76_yAqNKA_1wV7on2AwjQIk3tjkzsnk-Z4",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCttH3mlTJKKhQT\ny8RJBVIZDXKBdJSx9Y2VSFkmEK6ASMvE5W4VW7qX41PMg5VtGlajQCWxhugJwV1H\npEKvKqU0RsIZ14CTgTuVekRiruCnf0DN6KM6WuEJNsQxCTICchFC4sOt7kyjtCHU\nYUlF4faAFqhyMLzjJpdss+QQngSEjH1s+alr+RySBSA5puEF0VFMpRMEjFwi4Jno\n0yPNFCOiYvz+tRIK6OhSdtfhMKAmmt+62/Op3v9O2GAhFVJO1JNl0Odxlf/rs3nB\nv+37ds6pz+7LDOzGSufAAoigMDRRn1a93AT2Q4tlz9kYe5RncKhnA6xPgomNZfIf\nVMwjOGF1AgMBAAECggEAARbfr0GCRjrLU3B0s6yH3kZaUHuFrzQGBkik3ns+TOmn\n9X0m2pVvryIq1V6B4mRG5NEzK1DYRa9jwV5DWMvgq1pCP109ni8yS3av1RqZqBNB\nOclatLP7M06XnmMbYC6M8ylu5rlW27P2fll51ylanWUG+2hY1ufYDUN3i68iAh7I\nwsMIPdf1+TYGi8IdHg2hHpBEMecMU/bhRaDXM3DOIMWYZ8l6k7uxrsDKv7Uhov+n\nZIfv+AjWqz7N+yvIR6GM0HqG+iv4B+JAkrcYPbchl44IPbbo4HKRjEY2+LkOl/5I\nwM46pO6eAuUNXhXdw0y3WRqOdc/o+c7PD/JeZWTgpQKBgQDyu7hsKWapbpeZ5Kqm\nBxr+yec0/xcfXxkhcc9Ehdd8+v/QGKza8vRkCUp9elpccMpMW5rLdPVZpP2UovXF\nrSjDKl9m2H2OKzzBRO8uT/a//XSzywlQjcmSaYmu0stBhf7tSoDjlUjFqvYjdxDO\n1u2cqZiymVyC3R/bQv1RPm08zwKBgQC3MvDjBWeV6YBcu/iSF8jmokSOz71YyUFn\n9CReKzQzEDO6DVhWulsIDlPne3dfqI2vST+ZljNUo1S8jxaMX8C2rTdKmIGVWOQP\ntbQftdtti5QJ83Z5HhE4KSZ6irq1GY/AmL1DiU2mRHch7KOMKDGrTps1Rwv/Y8t2\n72Ha9q+2ewKBgQC88JwELVHJDtmYo6Kla6B6tSRwXyNbewWvv8wLVXc/xIy9KYfb\nQgQznfvKoiOWEwGU4DUkq5yTM9djDFnsjfXNvLzX7CoHMOawtfzLetjh5uMhVCii\n+Erv2ZCfcVtfXHLrt/ONstUbcBD52CNQLYJ1UJoYY0HcZ0z1ujY+OC6FhwKBgCpx\n0v3GMsm439SceGrgt9s3nUq5NtVrS4waNJLcz6tFBbcFgIIXix/Csg3fvTichLcn\n8WRUOHBTpz5IqKC9TpkEaNsPmnZPsgcxwhnWuJAY1qO3lKtbHAI3BoM9wSRUV8n3\nmWIcXbE4C6IAgaPnbBqUi8E8RLtXE7zqmXFx1iQhAoGAS0OlIv/Nr7K0cWdQouAT\nO7ohkipKJBKgD6rdU3e2lOtl3vheW7hxipe7id5hMwibeplbxEco1voxrn2qHiIv\n+g/WoSEOUkukzG7hyrlHL26L5Lpr3pTiTLpw3G3edyZqgR+mOIweEmQtQxIyOuyP\nepWwT2VHi5YaRSYwKgklrPk=\n-----END PRIVATE KEY-----",
  "default_bot_id": "7483341266928468008",
  "space_id": "7477025304549294118"
}
```

## 使用方法

### 1. 智能体对话
使用普通的模型名称调用智能体：

```bash
curl -X POST "http://localhost:3000/v1/chat/completions" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4",
    "messages": [
      {"role": "user", "content": "你好"}
    ],
    "stream": false
  }'
```

### 2. 同步工作流调用
使用 `workflow:` 前缀加工作流 ID：

```bash
curl -X POST "http://localhost:3000/v1/chat/completions" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "workflow:your-workflow-id",
    "messages": [
      {"role": "user", "content": "处理这个请求"}
    ],
    "stream": false
  }'
```

### 3. 异步工作流调用
使用 `workflow-async:` 前缀加工作流 ID，适用于长时间运行的工作流：

```bash
curl -X POST "http://localhost:3000/v1/chat/completions" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "workflow-async:your-long-running-workflow-id",
    "messages": [
      {"role": "user", "content": "处理复杂的长时间任务"}
    ],
    "stream": false
  }'
```

### 4. 流式响应
设置 `stream: true` 启用流式响应：

```bash
curl -X POST "http://localhost:3000/v1/chat/completions" \
  -H "Authorization: Bearer your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4",
    "messages": [
      {"role": "user", "content": "你好"}
    ],
    "stream": true
  }'
```

## 技术实现

### 1. JWT Token 生成
- 使用 RSA-256 算法签名
- Token 有效期为 1 小时
- 包含必要的 claims：iss, aud, iat, exp, jti

### 2. OAuth 流程
1. 生成 JWT token
2. 发送 OAuth 请求获取 access token
3. 缓存 access token（提前 5 分钟过期）
4. 使用 access token 调用 Coze API

### 3. 请求处理
- **智能体**：发送创建消息请求 → 轮询状态 → 获取消息详情
- **同步工作流**：发送请求 → 自动检测执行模式 → 返回结果
- **异步工作流**：发送请求 → 轮询执行状态 → 获取最终结果

### 4. 异步工作流轮询机制
- **轮询间隔**：2秒
- **最大轮询次数**：150次（总计约5分钟）
- **状态检测**：running → completed/failed/canceled
- **自动重试**：网络错误时自动重试
- **超时处理**：超过最大轮询次数时返回超时错误

## 注意事项

### 1. 安全性
- 私钥信息敏感，请妥善保管
- 建议定期轮换密钥
- 生产环境中使用 HTTPS

### 2. 限制
- JWT token 有效期为 1 小时
- Access token 会自动缓存和刷新
- 工作流 ID 需要正确配置

### 3. 错误处理
- 认证失败会返回相应错误信息
- 网络错误会自动重试
- 配置错误会在启动时报告

## 故障排除

### 1. 认证失败
- 检查 client_id 是否正确
- 验证私钥格式是否正确
- 确认公钥 ID 是否匹配

### 2. 请求失败
- 检查 Base URL 是否正确
- 验证智能体 ID 或工作流 ID
- 查看错误日志获取详细信息

### 3. 配置问题
- 确保 JSON 格式正确
- 检查必需字段是否完整
- 验证私钥是否包含正确的换行符

## 参考资料

- [Coze OAuth JWT 文档](https://www.coze.cn/open/docs/developer_guides/oauth_jwt)
- [Coze Go SDK 示例](https://github.com/coze-dev/coze-go/blob/main/examples/auth/jwt_oauth/main.go)
- [JWT 规范](https://tools.ietf.org/html/rfc7519)
